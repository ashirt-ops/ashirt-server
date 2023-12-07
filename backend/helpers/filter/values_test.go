package filter_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/helpers/filter"
	"github.com/stretchr/testify/require"
)

func TestValues(t *testing.T) {
	val1 := "plain"
	val2 := "modified"
	values := filter.Values{
		filter.Val(val1),
		filter.Value{Value: val2, Modifier: filter.Not},
	}

	require.Equal(t, []string{val1, val2}, values.Values())
}

func TestSplitValuesByModifier(t *testing.T) {
	val1 := "plain"
	val2 := "modified"
	values := filter.Values{
		filter.Val(val1),
		filter.Value{Value: val2, Modifier: filter.Not},
	}

	require.Equal(t, map[filter.FilterModifier][]string{
		filter.Normal: {val1},
		filter.Not:    {val2},
	}, values.SplitByModifier())
}

func TestSplitValues(t *testing.T) {
	val1 := "plain"
	val2 := "modified"
	values := filter.Values{
		filter.Val(val1),
		filter.Value{Value: val2, Modifier: filter.Not},
	}

	require.Equal(t, map[string][]string{
		"short": {val1},
		"long":  {val2},
	}, values.SplitValues(func(v filter.Value) string {
		if len(v.Value) < 6 {
			return "short"
		}
		return "long"
	}))
}

func TestValue(t *testing.T) {
	val1 := "plain"
	val2 := "modified"
	values := filter.Values{
		filter.Val(val1),
		filter.Value{Value: val2, Modifier: filter.Not},
	}

	require.Equal(t, val2, values.Value(1))
}
