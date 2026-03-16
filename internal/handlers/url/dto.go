package url

type GetRequest struct {
	ShortenedURL string `json:"url"`
}

type GetResponse struct {
	SourceURL string `json:"result"`
}
