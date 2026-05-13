package pagination

type Page struct {
	Page     int
	PageSize int
	Offset   int
}

func New(page, pageSize int) Page {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return Page{Page: page, PageSize: pageSize, Offset: (page - 1) * pageSize}
}
