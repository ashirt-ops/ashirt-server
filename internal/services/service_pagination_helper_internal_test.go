package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaginationSetMaxItems(t *testing.T) {
	p := Pagination{
		PageSize: 10000,
		Page:     1,
	}
	p.SetMaxItems(12)
	require.Equal(t, int64(12), p.maxPageSize)
}

func TestPaginationConstrain(t *testing.T) {
	p := Pagination{
		PageSize: 10000,
		Page:     0,
	}
	p.SetMaxItems(12)

	p.constrain()

	require.Equal(t, int64(12), p.PageSize)
	require.Equal(t, int64(1), p.Page)

	p = Pagination{
		PageSize: 10000,
		Page:     0,
	}

	p.constrain()
	require.Equal(t, int64(defaultMaxPageSize), p.PageSize)
}
