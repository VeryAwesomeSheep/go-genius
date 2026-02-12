package genius

type ReferentsService Service

type Referent struct {
	AnnotatorID    int          `json:"annotator_id"`
	AnnotatorLogin string       `json:"annotator_login"`
	APIPath        string       `json:"api_path"`
	Classification string       `json:"classification"`
	Fragment       string       `json:"fragment"`
	ID             int          `json:"id"`
	IsDescription  bool         `json:"is_description"`
	SongID         int          `json:"song_id"`
	URL            string       `json:"url"`
	Annotatable    *Annotatable `json:"annotatable"`
}

type Annotatable struct {
	APIPath  string `json:"api_path"`
	Context  string `json:"context"`
	ID       int    `json:"id"`
	ImageURL string `json:"image_url"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type ReferentsOptions struct {
	CreatedByID int `url:"created_by_id,omitempty"`
	SongID      int `url:"song_id,omitempty"`
	WebPageID   int `url:"web_page_id,omitempty"`
	TextFormatOptions
	PagingOptions
}
