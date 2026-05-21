package list

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// Request holds the pagination query parameters for the list endpoint.
type Request struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func (r *Request) normalize() {
	if r.Page < 1 {
		r.Page = defaultPage
	}
	if r.PageSize < 1 {
		r.PageSize = defaultPageSize
	}
	if r.PageSize > maxPageSize {
		r.PageSize = maxPageSize
	}
}
