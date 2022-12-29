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

type URLBatchInput struct {
	CorrelID string `json:"correlation_id"`
	URL      string `json:"original_url"`
}

type URLBatchResponse struct {
	CorrelID string `json:"correlation_id"`
	URL      string `json:"short_url"`
}
