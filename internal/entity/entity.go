package entity

type ShortURL struct {
	ID     string
	Origin string
}

type URLInput struct {
	URL string `json:"url"`
}

type URLResponse struct {
	Result string `json:"result"`
}
