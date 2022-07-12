// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/theparanoids/ashirt-server/backend/helpers"
)

func TestChanToMap(t *testing.T) {
	values := []int{4, 8, 15, 16, 23, 42}
	ch := make(chan int, len(values))

	for _, v := range values {
		ch <- v
	}

	result := helpers.ChanToSlice(&ch)
	require.Equal(t, values, result)
}
