// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package fuzzyfilefinder

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/agnivade/levenshtein"
)

var (
	// ErrEmptyDirectory is returned when the directory contains no files
	ErrEmptyDirectory = errors.New("no files found in directory")
	// ErrAmbiguousDistance is returned when multiple files have the same
	// levenshtein distance
	ErrAmbiguousDistance = errors.New("multiple files found within distance place; remove unused files from the working directory")
	// ErrNoFiles is returned when no files are found within the specified
	// levenshtein distance
	ErrNoFiles = errors.New("no files found within levenshtein distance defined")
)

//this struct is used to store levenshtiencompairsons

type levenshteinValues struct {
	// file name is the name of the file you are compairing
	filename string
	//compairson in the original file name of the file your trying to find
	compairson string
	//this is the levenstien distance
	distance int
}

// LvSorter sorts LevenshteinValues by distance.
type LvSorter []levenshteinValues

//sort array of lv values to find the one with the least distance
func (a LvSorter) Len() int           { return len(a) }
func (a LvSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LvSorter) Less(i, j int) bool { return a[i].distance < a[j].distance }

// Exists checks whether a file exists on the filesystem
func Exists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}

// FuzzyFinder identifies files that exist on the filesystem who's names are
// within a specific levenshtein distance
func FuzzyFinder(path string, lvDistance int) (string, error) {
	var lv = make([]levenshteinValues, 0)

	if Exists(path) {
		return path, nil
	}

	missingFile := filepath.Base(path)
	dir := filepath.Dir(path)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", ErrEmptyDirectory
	}

	//check distance on all files in directory
	for _, f := range files {
		distance := levenshtein.ComputeDistance(f.Name(), missingFile)
		flv := levenshteinValues{f.Name(), missingFile, distance}
		lv = append(lv, flv)
	}

	sort.Sort(LvSorter(lv))
	if len(lv) > 1 && lv[0].distance == lv[1].distance {
		return "", ErrAmbiguousDistance
	}

	if lvDistance > lv[0].distance {
		return filepath.Join(dir, lv[0].filename), nil
	}

	return "", ErrNoFiles
}
