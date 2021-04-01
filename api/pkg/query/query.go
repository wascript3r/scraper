package query

type QueryRes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type GetAllRes struct {
	Queries []*QueryRes `json:"queries"`
}
