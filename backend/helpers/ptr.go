package helpers

// StringPtr converts a string into a *string
func StringPtr(s string) *string {
	return &s
}

// I64Ptr converts a int64 into an *int64
func I64Ptr(i int64) *int64 {
	return &i
}
