package server

type BookRequestQuery struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Isbn      string `json:"isbn"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
