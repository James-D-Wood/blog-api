package model

type BlogPost struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Contents    string `json:"contents"`
	AuthorID    string `json:"author_id"`
	CreatedTS   string `json:"created_ts"`
	PublishedTS string `json:"published_ts"`
	UpdatedTS   string `json:"updated_ts"`
}
