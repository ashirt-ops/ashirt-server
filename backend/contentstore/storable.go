// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

// Storable represents content in a yet-to-be-uploaded state. The basic workflow when using this is
// to create a new instance, then immediately call ProcessPreviewAndUpload, which will both generate
// preview content (e.g. a resized image), and upload both the provided content, plus the preview
// to the provided Store
type Storable interface {
	ProcessPreviewAndUpload(Store) (ContentKeys, error)
}
