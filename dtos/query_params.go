package dtos

type QueryParams struct {
	Page    int64
	PerPage int64
	OrderBy string
	Filter  string
}
