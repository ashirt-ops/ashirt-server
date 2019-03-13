package isthere

// No "asserts" that an error is nil. Returns true if the error is nil, false otherwise
// Why not just err == nil? Most ifs use the latter : err != nil. This can make it
// challenging when reviewing if the check should be err == nil vs err != nil. This solves
// that problem by instead chaning the format: No(err)
func No(e error) bool {
	return e == nil
}
