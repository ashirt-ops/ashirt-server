package filter

type FilterModifier int64

const (
	Normal FilterModifier = 0
	Not    FilterModifier = 1 << (iota)
)
