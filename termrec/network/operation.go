package network

// Operation is linked to the AShirt operation structure retrieved from GET /web/operations
type Operation struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	NumUsers int    `json:"numUsers"`
	Status   int    `json:"status"`

	// ID is only used in list operations for the API since the screenshot client still expects int64 ids.
	// Once the screenshot client is updated this line can be removed
	ID int64 `json:"id"`
}
