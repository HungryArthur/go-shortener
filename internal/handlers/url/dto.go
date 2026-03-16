package url

type GetRequest struct {
	ShortenedURL string `json:"url"`
}

type GetResponse struct {
	SourceURL string `json:"result"`
}

type CreateRequest struct {
	SourceURL string `json:"url"`
}

type CreateResponse struct {
	ShortenedURL string `json:"result"`
}
