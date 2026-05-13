package response

type Envelope[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
}

type PageMeta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type Page[T any] struct {
	Data T        `json:"data"`
	Meta PageMeta `json:"meta"`
}
