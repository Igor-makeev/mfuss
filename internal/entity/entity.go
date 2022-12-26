package entity

type ShortURL struct {
	ID        string `json:"-"`
	ResultURL string `json:"short_url"`
	Origin    string `json:"original_url"`
	UserID    string `json:"-"`
}

type URLInput struct {
	URL string `json:"url"`
}

type URLResponse struct {
	Result string `json:"result"`
}
