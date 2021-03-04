package query

type QueryRes struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type GetAllRes struct {
	Queries []*QueryRes `json:"queries"`
}
