package processors

import (
	"archive/zip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/logging"
	"github.com/theparanoids/ashirt/backend/models"
	"golang.org/x/sync/errgroup"

	sq "github.com/Masterminds/squirrel"
)

// ExportProcessor is a struct that can start a long-running process to watch the exports_queue table
// and properly export those items to the archive store
type ExportProcessor struct {
	db            *database.Connection
	contentStore  contentstore.Store
	archiveStore  contentstore.Store
	stopChan      chan bool
	isStoppedChan chan bool
	logger        logging.Logger
}

var exportTemplates = "/app/backend/archive_templates/"

// NewExportProcessor generates a new ExportProcessor, configured with the passed parameters
func NewExportProcessor(db *database.Connection, contentStore, archiveStore contentstore.Store, isStoppedChan chan bool, logger logging.Logger) *ExportProcessor {
	ep := ExportProcessor{
		db:            db,
		contentStore:  contentStore,
		archiveStore:  archiveStore,
		stopChan:      make(chan bool),
		isStoppedChan: isStoppedChan,
		logger:        logging.With(logger, "process", "export"),
	}
	return &ep
}

// Start kicks off the processor. This method will return immediately, but also start a background
// go routine that will do the real work
func (ep *ExportProcessor) Start() {
	go ep.process()
}

// Stop kills the processor. Note that this is not immediate, and instead will stop at the next
// check event (either after a complete export, or, if not processing, after it's check-delay interval)
func (ep *ExportProcessor) Stop() {
	ep.stopChan <- true
}

func (ep *ExportProcessor) process() {
	keepRunning := true
	for keepRunning {
		select {
		case <-ep.stopChan:
			keepRunning = false
		default:
			processed := ep.checkQueue()
			if !processed {
				ep.logger.Log("msg", "Waiting for stuff")
				time.Sleep(time.Minute)
			}
		}
	}
	ep.isStoppedChan <- true
}

func (ep *ExportProcessor) checkQueue() bool {
	var exportItem models.ExportQueueItem
	var updateExportBase sq.UpdateBuilder
	err := ep.db.WithTx(context.Background(), func(tx *database.Transactable) { //using a transaction to lock the row while finding an operation
		tx.Get(&exportItem, sq.Select("id", "operation_id").From("exports_queue").Where(
			sq.Eq{"status": models.ExportStatusPending},
			sq.Gt{"created_at": time.Minute},
		).OrderBy("id ASC").Limit(1).Suffix("FOR UPDATE"))

		updateExportBase = sq.Update("exports_queue").Where(sq.Eq{"id": exportItem.ID})

		tx.Update(updateExportBase.Set("status", models.ExportStatusInProgress))
	})
	if err != nil {
		if err != sql.ErrNoRows {
			ep.logger.Log("queueCheck", "failed", "error", err.Error())
		}
		return false
	}

	exportLog := logging.With(logging.With(ep.logger, "exportID", exportItem.ID), "operationID", exportItem.OperationID)
	exportLog.Log("exportUpdate", "Starting")

	err = ep.doExport(exportTemplates, exportItem.OperationID, updateExportBase)
	var newStatus models.ExportStatus
	var notes *string
	if err != nil {
		exportLog.Log("exportUpdate", "Failed", "error", err.Error())
		newStatus = models.ExportStatusError
		temp := err.Error()
		notes = &temp
	} else {
		exportLog.Log("exportUpdate", "Complete")
		newStatus = models.ExportStatusComplete
	}

	err = ep.db.Update(updateExportBase.Set("status", newStatus).Set("notes", notes))
	if err != nil {
		exportLog.Log("msg", "Unable to update export with new status", "error", err.Error(), "status", newStatus, "notes", notes)
	}
	return true
}

// doExport gathers all data from the database and content store to provide a snapshot of an existing
// operation at the time this gets called. This is a long process, and could fail at various places along the
// way.
func (ep *ExportProcessor) doExport(staticFileDir string, operationID int64, updateExportQuery sq.UpdateBuilder) error {

	export, err := extractOperationData(ep.db, operationID)
	if err != nil {
		return err
	}

	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()
	zipWriter := zip.NewWriter(pw)
	defer func() {
		zipWriter.Flush()
		zipWriter.Close()
	}()

	exportName := fmt.Sprintf("%v_%v.zip", export.Slug, time.Now().Unix())
	ep.db.Update(updateExportQuery.Set("export_name", exportName))

	var g errgroup.Group
	g.Go(func() error { return ep.archiveStore.UploadWithName(exportName, pr) })

	archiveRoot := export.Name + "/"
	assetsRoot := archiveRoot + "media/"

	jsonEvidencePrefix := []byte("evidenceJsonp(")
	jsonEvidencePostfix := []byte(")\n")

	// migrate content
	for _, key := range MapEvidenceToContentKeys(export.Evidence) {
		if key.ContentType == "image" {
			if err = copyContentStoreFile(zipWriter, assetsRoot, key.Key, ep.contentStore); err != nil {
				return backend.UploadErr(err)
			}
		} else if key.ContentType == "terminal-recording" {
			if err = copyContentStoreFileAsString(zipWriter, assetsRoot, key.Key, ep.contentStore, jsonEvidencePrefix, jsonEvidencePostfix); err != nil {
				return backend.UploadErr(err)
			}
		} else {
			if err = copyContentStoreFileWithFixes(zipWriter, assetsRoot, key.Key, ep.contentStore, jsonEvidencePrefix, jsonEvidencePostfix); err != nil {
				return backend.UploadErr(err)
			}
		}
	}

	if err = writeOperationJSON(zipWriter, archiveRoot, *export); err != nil {
		return backend.UploadErr(err)
	}

	relativeRoot := filepath.Clean(staticFileDir + "/")
	err = filepath.Walk(relativeRoot, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return fmt.Errorf("path (%v) does not exist", path)
		}
		if !info.IsDir() {
			if info.Name() == ".gitkeep" || info.Name() == ".gitignore" {
				return nil // ignore this file
			}

			relativePath := StripPathPrefix(relativeRoot, path)
			if strings.HasPrefix(relativePath, "/") {
				relativePath = relativePath[1:]
			}
			return copyFileToZip(zipWriter, path, archiveRoot+relativePath)
		}
		return nil
	})
	if err != nil {
		return backend.UploadErr(err)
	}

	return firstNonNilError(zipWriter.Close(), pw.Close(), g.Wait(), pr.Close())
}

func firstNonNilError(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

// copyContentStoreFile copies the content of a given key into the same file name locally
// (for export purposes)
func copyContentStoreFile(zipWriter Creator, dir, key string, contentStore contentstore.Store) error {
	return copyContentStoreFileWithFixes(zipWriter, dir, key, contentStore, []byte{}, []byte{})
}

// copyContentStoreFileWithFixes is the same as copyContentStoreFile, but with a prefix and postfix applied.
func copyContentStoreFileWithFixes(zipWriter Creator, dir, key string, contentStore contentstore.Store, prefix, postfix []byte) error {
	contentReader, err := contentStore.Read(key)
	if err != nil {
		return err
	}

	return copyStreamToZipWithFixes(zipWriter, contentReader, dir+key, prefix, postfix)
}

func copyContentStoreFileAsString(zipWriter Creator, dir, key string, contentStore contentstore.Store, prefix, postfix []byte) error {
	contentReader, err := contentStore.Read(key)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(contentReader)
	if err != nil {
		return err
	}
	content, err = json.Marshal(string(content))
	if err != nil {
		return err
	}

	to, err := zipWriter.Create(dir + key)
	if err != nil {
		return err
	}

	_, err = to.Write(prefix)
	if err != nil {
		return err
	}
	_, err = to.Write(content)
	if err != nil {
		return err
	}
	_, err = to.Write(postfix)
	return err
}

// copyFileToZip opens a local file (located at src), creates a file in the zipfile at dst, and copies
// all bytes from src to dst
func copyFileToZip(zipWriter Creator, src, dst string) error {
	from, err := os.Open(src) // opting to re-open every time to avoid unnecessary memory use
	if err != nil {
		return err
	}
	return firstNonNilError(copyStreamToZip(zipWriter, from, dst), from.Close())
}

// copyStreamToZip creates a new file (at dst) and moves the data in src into the newly created file
func copyStreamToZip(zipWriter Creator, src io.Reader, dst string) error {
	return copyStreamToZipWithFixes(zipWriter, src, dst, []byte{}, []byte{})
}

// copyStreamToZipWithFixes is the same as copyStreamToZip, but with a prefix and postfix applied.
func copyStreamToZipWithFixes(zipWriter Creator, src io.Reader, dst string, prefix, postfix []byte) error {
	to, err := zipWriter.Create(dst)
	if err != nil {
		return err
	}

	_, err = to.Write(prefix)
	if err != nil {
		return err
	}
	_, err = io.Copy(to, src)
	if err != nil {
		return err
	}
	_, err = to.Write(postfix)
	return err
}

// extractOperationData retrieves, in parallel, all of the direct and indirect data for an operation
// from the database
func extractOperationData(db *database.Connection, operationID int64) (*models.OperationExport, error) {
	var export models.OperationExport
	var g errgroup.Group
	selectStarByOperationID := sq.Select("*").Where(sq.Eq{"operation_id": operationID})

	g.Go(func() error {
		return db.Get(&export,
			sq.Select("id", "slug", "name", "status", "created_at", "updated_at").
				From("operations").Where(sq.Eq{"id": operationID}))
	})
	g.Go(func() error { return db.Select(&export.Queries, selectStarByOperationID.From("queries")) })
	g.Go(func() error { return db.Select(&export.Evidence, selectStarByOperationID.From("evidence")) })
	g.Go(func() error { return db.Select(&export.Findings, selectStarByOperationID.From("findings")) })
	g.Go(func() error { return db.Select(&export.Tags, selectStarByOperationID.From("tags")) })

	if err := g.Wait(); err != nil {
		return nil, backend.DatabaseErr(err)
	}

	evidenceIDs := make([]int64, len(export.Evidence))
	relevantUserIDs := make([]int64, len(export.Evidence)) // this may have duplicate entries, but sql will de-dupe this when querying

	for i, e := range export.Evidence {
		evidenceIDs[i] = e.ID
		relevantUserIDs[i] = e.OperatorID
	}

	selectStarByEvidenceIDs := sq.Select("*").Where(sq.Eq{"evidence_id": evidenceIDs})

	g.Go(func() error {
		return db.Select(&export.EvidenceToFinding, selectStarByEvidenceIDs.From("evidence_finding_map"))
	})
	g.Go(func() error {
		return db.Select(&export.TagsToEvidence, selectStarByEvidenceIDs.From("tag_evidence_map"))
	})
	g.Go(func() error {
		return db.Select(&export.Users, sq.Select("*").From("users").Where(sq.Eq{"id": relevantUserIDs}))
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &export, nil
}

// MapEvidenceToContentKeys retrieves all non-empty, unique values of Evidence.ThumbImageKey and
// Evidence.FullImageKey.
// returns a slice of these keys, plus the contentType of the evidence
func MapEvidenceToContentKeys(evidence []models.Evidence) []EvidenceKeyContentTypePair {
	evidenceContentKeys := make([]EvidenceKeyContentTypePair, 0, 2*len(evidence))
	for _, e := range evidence {
		if e.FullImageKey != "" {
			evidenceContentKeys = append(evidenceContentKeys, EvidenceKeyContentTypePair{e.FullImageKey, e.ContentType})
		}
		if e.FullImageKey != e.ThumbImageKey && e.ThumbImageKey != "" {
			evidenceContentKeys = append(evidenceContentKeys, EvidenceKeyContentTypePair{e.ThumbImageKey, e.ContentType})
		}
	}
	return evidenceContentKeys
}

// writeOperationJSON serializes the OperationExport and writes this to data.json file
func writeOperationJSON(zipWriter Creator, root string, data models.OperationExport) error {
	encodedData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	contents := "archiveJsonp(" + string(encodedData) + ")\n"
	w, err := zipWriter.Create(root + "data.json")
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(contents))
	return err
}

// Creator is a simple interface for abstracting relevant portions of zip.Writer.
type Creator interface {
	Create(string) (io.Writer, error)
}

// StripPathPrefix removes a leading prefix from a given string
//
// Note: This is specifically for file paths, and this does not check
// if the provided fullpath contains the prefix, it just removes it.
func StripPathPrefix(prefix, fullpath string) string {
	cleanedPath := filepath.Clean(fullpath)
	relativeToRoot := cleanedPath[len(prefix):]
	return relativeToRoot
}

// EvidenceKeyContentTypePair represents Evidence as a key/contentType pairing
type EvidenceKeyContentTypePair struct {
	Key         string
	ContentType string
}
