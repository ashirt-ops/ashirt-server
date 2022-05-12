// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package filter

type FilterModifier int64

const (
	Normal FilterModifier = 0
	Not    FilterModifier = 1 << (iota)
)
