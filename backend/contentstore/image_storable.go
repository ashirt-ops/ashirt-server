// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

import (
	"bytes"
	"image"

	// provides handling for decoding images
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
	"github.com/nfnt/resize"
	"golang.org/x/sync/errgroup"
)

type imageStorable struct {
	data io.Reader
}

// NewImage returns a Storable that deals specifically with images (png / jpg)
func NewImage(data io.Reader) Storable {
	return imageStorable{
		data: data,
	}
}

// ProcessPreviewAndUpload alters an image to fit within a specified size (currently hard coded to 500x500)
// and then uploads to the backing image store. The input image is expected to be in jpeg format,
// while the resized thumbnail will be in png format
func (is imageStorable) ProcessPreviewAndUpload(s Store) (ContentKeys, error) {
	var g errgroup.Group
	contentKeys := ContentKeys{}

	imageBytes, err := io.ReadAll(is.data)
	if err != nil {
		return contentKeys, err
	}

	g.Go(func() (err error) {
		contentKeys.Full, err = s.Upload(bytes.NewReader(imageBytes))
		return err
	})

	g.Go(func() (err error) {
		resized, err := resizeImage(500, 500, bytes.NewReader(imageBytes))
		if err != nil {
			return err
		}
		contentKeys.Thumbnail, err = s.Upload(resized)
		return err
	})

	return contentKeys, g.Wait()
}

func resizeImage(maxWidth uint, maxHeight uint, data io.Reader) (io.Reader, error) {
	img, _, err := image.Decode(data)
	if err != nil {
		return nil, err
	}

	size := img.Bounds().Size()
	newWidth := uint(0)
	newHeight := uint(0)
	if float32(maxWidth)/float32(size.X) < float32(maxHeight)/float32(size.Y) {
		newWidth = maxWidth
	} else {
		newHeight = maxHeight
	}

	m := resize.Resize(newWidth, newHeight, img, resize.NearestNeighbor)
	out := &bytes.Buffer{}
	err = png.Encode(out, m)
	return bytes.NewReader(out.Bytes()), err
}
