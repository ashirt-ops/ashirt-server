// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package processors_test

// import (
// 	"archive/zip"
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"testing"

// 	"github.com/theparanoids/ashirt/backend/contentstore"
// 	"github.com/theparanoids/ashirt/backend/models"
// 	"github.com/theparanoids/ashirt/backend/services"
// 	"github.com/stretchr/testify/require"
// )

// func TestExportOperation(t *testing.T) {
// 	db := initTest(t)
// 	HarryPotterSeedData.ApplyTo(t, db)
// 	contentStore, _ := contentstore.NewMemStore()
// 	archiveStore, _ := contentstore.NewMemStore()
// 	fileStore, _ := contentstore.NewMemStore() // helper to store contents of static files
// 	adminUser := UserDumbledore
// 	ctx := simpleFullContext(adminUser)

// 	must := muster(t)

// 	targetOp := OpChamberOfSecrets
// 	keysToPopulate := services.MapEvidenceToContentKeys(evidenceForOperation(t, HarryPotterSeedData, targetOp.ID))
// 	populateEvidenceContentIntoStore(keysToPopulate, contentStore)

// 	staticDirPath := "fixtures/"
// 	populateStaticContentIntoStore(t, staticDirPath, fileStore)

// 	must(services.ExportOperation(ctx, db, contentStore, archiveStore, staticDirPath, targetOp.Slug))

// 	//convert zip stream into reader
// 	archivedZipReader := must(archiveStore.Read(targetOp.Slug + ".zip")).(io.Reader)
// 	zipReader := unzip(t, archivedZipReader)

// 	allZipPrefix := targetOp.Name + "/"
// 	mediaPrefix := allZipPrefix + "media/"
// 	foundData := false
// 	foundMedia := make([]string, 0, len(keysToPopulate))
// 	foundStatic := make([]string, 0, len(fileStore.Keys()))
// 	for _, zippedFile := range zipReader.File {
// 		fileReader := must(zippedFile.Open()).(io.ReadCloser)
// 		actualFileBytes := must(ioutil.ReadAll(fileReader)).([]byte)

// 		switch {
// 		case strings.HasPrefix(zippedFile.Name, mediaPrefix):
// 			contentStoreKey := filepath.Base(zippedFile.Name)
// 			foundMedia = append(foundMedia, contentStoreKey)
// 			require.Equal(t, readFromStore(t, contentStoreKey, contentStore), actualFileBytes)
// 		case zippedFile.Name == allZipPrefix+"data.json":
// 			foundData = true
// 			contentAsString := string(actualFileBytes)
// 			expectedPrefix, expectedSuffix := "archiveJsonp(", ")\n"
// 			require.True(t, strings.HasPrefix(contentAsString, expectedPrefix))
// 			require.True(t, strings.HasSuffix(contentAsString, expectedSuffix))
// 			trimmed := contentAsString[len(expectedPrefix) : len(contentAsString)-len(expectedSuffix)]
// 			var export models.OperationExport
// 			require.NoError(t, json.Unmarshal([]byte(trimmed), &export))
// 			require.Equal(t, targetOp.ID, export.ID)
// 		default: // rest must be static files
// 			filestoreKey := services.StripPathPrefix(allZipPrefix, zippedFile.Name)
// 			foundStatic = append(foundStatic, filestoreKey)
// 			require.Equal(t, readFromStore(t, filestoreKey, fileStore), actualFileBytes)
// 		}
// 		require.NoError(t, fileReader.Close())
// 	}
// 	require.True(t, foundData)
// 	require.Equal(t, foundMedia, keysToPopulate)
// 	require.Equal(t, sortedStrings(foundStatic), sortedStrings(fileStore.Keys()))
// }

// func TestMapEvidenceToContentKeys(t *testing.T) {
// 	evidence := []models.Evidence{
// 		models.Evidence{FullImageKey: "Alpha", ThumbImageKey: "Beta"},
// 		models.Evidence{FullImageKey: "Gamma", ThumbImageKey: "Gamma"},
// 		models.Evidence{FullImageKey: "Tau", ThumbImageKey: "Rho"},
// 		models.Evidence{FullImageKey: "", ThumbImageKey: "Kappa"},
// 		models.Evidence{FullImageKey: "Delta", ThumbImageKey: ""},
// 		models.Evidence{FullImageKey: "", ThumbImageKey: ""},
// 	}
// 	require.Equal(t,
// 		[]string{"Alpha", "Beta", "Gamma", "Tau", "Rho", "Kappa", "Delta"},
// 		services.MapEvidenceToContentKeys(evidence),
// 	)
// }

// func TestStripPathPrefix(t *testing.T) {
// 	prefix := "root/path/"
// 	fullpath := "root/path/to/file"

// 	require.Equal(t, "to/file", services.StripPathPrefix(prefix, fullpath))
// }

// /////////// Test Helpers
// func readFromStore(t *testing.T, key string, store contentstore.Store) []byte {
// 	must := muster(t)
// 	contentReader := must(store.Read(key)).(io.Reader)
// 	return must(ioutil.ReadAll(contentReader)).([]byte)
// }

// func populateEvidenceContentIntoStore(keys []string, store contentstore.Store) {
// 	for _, key := range keys {
// 		content := bytes.NewBuffer([]byte(key + "--" + key)) // add some unique junk to verify content
// 		store.UploadWithName(key, content)
// 	}
// }

// func populateStaticContentIntoStore(t *testing.T, dir string, store contentstore.Store) {
// 	must := muster(t)
// 	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
// 		if !info.IsDir() {
// 			relativePath := services.StripPathPrefix(dir, path)
// 			f := must(os.Open(path)).(*os.File)
// 			store.UploadWithName(relativePath, f)
// 			require.NoError(t, f.Close())
// 		}
// 		return nil
// 	})
// }

// func unzip(t *testing.T, archivedZipReader io.Reader) *zip.Reader {
// 	must := muster(t)
// 	zipBytes := must(ioutil.ReadAll(archivedZipReader)).([]byte)
// 	return must(zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))).(*zip.Reader)
// }
