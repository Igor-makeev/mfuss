// Пакет с сущностями ошибок
package errorsentity

import "fmt"

// Тип кастомной ошибки на случай если ссылка уже существует в хранилище
type URLConflict struct {
	Str string
}

// Имплементация метода Error()
func (is URLConflict) Error() string {
	return fmt.Sprintf("error:  url: %v has already been shortened", is.Str)
}
