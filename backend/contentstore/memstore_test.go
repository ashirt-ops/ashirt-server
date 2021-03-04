// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
)

func TestMemstore(t *testing.T) {
	store, _ := contentstore.NewMemStore()

	content := []byte("Very innocent stuff")
	key, err := store.Upload(bytes.NewReader(content))
	require.NoError(t, err)
	require.NotEqual(t, "", key, "Key should be populated in response")

	reader, err := store.Read(key)
	require.NoError(t, err)

	b, _ := io.ReadAll(reader)

	require.Equal(t, content, b, "retrieved content should match uploaded content")
}

func TestMemstoreNoSuchKey(t *testing.T) {
	store, _ := contentstore.NewMemStore()

	_, err := store.Read("????")
	require.NotNil(t, err, "An error should be produced when a key is not found")
}
