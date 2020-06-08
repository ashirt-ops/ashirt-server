// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package processors_test

import (
	"bytes"

	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/processors"
)

func TestMapEvidenceToContentKeys(t *testing.T) {
	evidence := []models.Evidence{
		models.Evidence{ContentType: "1", FullImageKey: "Alpha", ThumbImageKey: "Beta"},
		models.Evidence{ContentType: "2", FullImageKey: "Gamma", ThumbImageKey: "Gamma"},
		models.Evidence{ContentType: "3", FullImageKey: "Tau", ThumbImageKey: "Rho"},
		models.Evidence{ContentType: "4", FullImageKey: "", ThumbImageKey: "Kappa"},
		models.Evidence{ContentType: "5", FullImageKey: "Delta", ThumbImageKey: ""},
		models.Evidence{ContentType: "6", FullImageKey: "", ThumbImageKey: ""},
	}
	require.Equal(t,
		[]processors.EvidenceKeyContentTypePair{
			{"Alpha", "1"}, 
			{"Beta", "1"}, 
			{"Gamma", "2"}, 
			{"Tau", "3"},
			{"Rho", "3"}, 
			{"Kappa", "4"},
			{"Delta", "5"},
		},
		processors.MapEvidenceToContentKeys(evidence),
	)
}

func TestStripPathPrefix(t *testing.T) {
	prefix := "root/path/"
	fullpath := "root/path/to/file"

	require.Equal(t, "to/file", processors.StripPathPrefix(prefix, fullpath))
}

/////////// Test Helpers
func populateEvidenceContentIntoStore(keys []string, store contentstore.Store) {
	for _, key := range keys {
		content := bytes.NewBuffer([]byte(key + "--" + key)) // add some unique junk to verify content
		store.UploadWithName(key, content)
	}
}
