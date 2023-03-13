// Пакет c сущностями
package entity

// Тип сокращщенной URL
type ShortURL struct {
	ID        string `json:"-"`
	ResultURL string `json:"short_url"`
	Origin    string `json:"original_url"`
	UserID    string `json:"-"`
	IsDeleted bool   `json:"-"`
}

// Тип входещей URL
type URLInput struct {
	URL string `json:"url"`
}

// Тип ответа URL
type URLResponse struct {
	Result string `json:"result"`
}

// Тип входящей пачки URL
type URLBatchInput struct {
	CorrelID string `json:"correlation_id"`
	URL      string `json:"original_url"`
}

// Тип ответа пачки URL
type URLBatchResponse struct {
	CorrelID string `json:"correlation_id"`
	URL      string `json:"short_url"`
}

// Функция проставляющая флаг удаления
func (s *ShortURL) SetDeleteFlag() {
	s.IsDeleted = true
}
