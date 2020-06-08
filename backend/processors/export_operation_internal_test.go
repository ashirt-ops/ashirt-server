// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package processors

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/models"
)

var internalClock = clockwork.NewFakeClock()

func TestFirstNonNilError(t *testing.T) {
	var dummyError1, dummyError2 error
	realError1 := errors.New("Red")
	realError2 := errors.New("Blue")
	require.Equal(t, realError1, firstNonNilError(dummyError1, dummyError2, realError1, realError2))
	require.Equal(t, realError1, firstNonNilError(dummyError1, realError1, dummyError2, realError2))
	require.Nil(t, firstNonNilError(dummyError1, dummyError2))
}

func TestWriteOperationJSON(t *testing.T) {
	testData := models.OperationExport{
		Operation: models.Operation{
			Name: "SomeName",
		},
		Tags: []models.Tag{
			models.Tag{
				ID:          1,
				OperationID: 2,
				Name:        "Tag1",
				ColorName:   "Red",
				CreatedAt:   internalClock.Now(),
				UpdatedAt:   nil,
			},
		},
	}
	creator := newSimpleCreator("errorPath/data.json")

	err := writeOperationJSON(&creator, "errorPath/", testData)
	require.Error(t, err)

	root := "fortesting/"
	err = writeOperationJSON(&creator, root, testData)
	require.NoError(t, err)
	require.Equal(t, root+"data.json", creator.StreamName)

	encodedData, _ := json.Marshal(testData)
	require.Equal(t, "archiveJsonp("+string(encodedData)+")\n", creator.Stream.String())
}

func TestCopyStreamToZip(t *testing.T) {
	creator := newSimpleCreator()
	content := "one two three"
	filename := "newfile"
	buf := bytes.NewBuffer([]byte(content))
	err := copyStreamToZip(&creator, buf, filename)
	require.NoError(t, err)
	require.Equal(t, filename, creator.StreamName)
	require.Equal(t, content, creator.Stream.String())
}

func TestCopyContentStoreFile(t *testing.T) {
	memStore, _ := contentstore.NewMemStore()
	content := "one two three"
	key, err := memStore.Upload(bytes.NewBuffer([]byte(content)))
	require.NoError(t, err)

	creator := newSimpleCreator()
	dir := "junk/"
	err = copyContentStoreFile(&creator, dir, key, memStore)
	require.NoError(t, err)
	require.Equal(t, content, creator.Stream.String())
	require.Equal(t, dir+key, creator.StreamName)
}

func TestCopyContentStoreFileAsString(t *testing.T) {
	memStore, _ := contentstore.NewMemStore()
	content := `one "two" three`
	prefix := []byte("zero")
	postfix := []byte("four")
	key, err := memStore.Upload(bytes.NewBuffer([]byte(content)))
	require.NoError(t, err)

	creator := newSimpleCreator()
	dir := "junk/"
	err = copyContentStoreFileAsString(&creator, dir, key, memStore, prefix, postfix)
	require.NoError(t, err)

	rawContent := creator.Stream.String()
	decodedPrefix := rawContent[:len(prefix)]
	decodedPostfix := rawContent[len(rawContent)-len(postfix):]
	jsonContent := rawContent[len(prefix) : len(rawContent)-len(postfix)]
	require.Equal(t, string(prefix), decodedPrefix)
	require.Equal(t, string(postfix), decodedPostfix)
	var unwrappedContent string
	require.NoError(t, json.Unmarshal([]byte(jsonContent), &unwrappedContent))
	require.Equal(t, content, unwrappedContent)
	require.Equal(t, dir+key, creator.StreamName)
}

// SimpleCreator provides a single, small, reusable buffer to testing synchronous bufer writing
type SimpleCreator struct {
	StreamName    string
	Stream        *bytes.Buffer
	NameBlacklist []string
}

func newSimpleCreator(nameBlacklist ...string) SimpleCreator {
	return SimpleCreator{
		Stream:        bytes.NewBuffer(make([]byte, 32*1024)),
		NameBlacklist: nameBlacklist,
	}
}

func (s *SimpleCreator) IsInBlacklist(name string) bool {
	for _, blacklistedName := range s.NameBlacklist {
		if name == blacklistedName {
			return true
		}
	}
	return false
}

func (s *SimpleCreator) Create(dst string) (io.Writer, error) {
	s.StreamName = ""
	s.Stream.Reset()
	if s.IsInBlacklist(dst) {
		return nil, errors.New("SimpleCreator trying to Create on blacklisted name")
	}
	s.StreamName = dst

	return s.Stream, nil
}

func TestCreator(t *testing.T) {
	clearPath, errorPath, altErrorPath := "clearPath", "errorPath", "altErrorPath"
	testMessage := "Okay!"
	creator := newSimpleCreator(errorPath, altErrorPath)

	require.True(t, creator.IsInBlacklist(errorPath))
	require.True(t, creator.IsInBlacklist(altErrorPath))
	require.False(t, creator.IsInBlacklist(clearPath))

	w, err := creator.Create(errorPath)
	require.Error(t, err)
	require.Nil(t, w)
	require.Equal(t, "", creator.StreamName)

	w, err = creator.Create(clearPath)
	require.NoError(t, err)
	require.NotNil(t, w)

	written, err := w.Write([]byte(testMessage))
	require.NoError(t, err)
	require.Equal(t, len(testMessage), written)

	require.Equal(t, clearPath, creator.StreamName)
	require.Equal(t, w, creator.Stream)
	require.Equal(t, testMessage, creator.Stream.String())
}
