package errorsentity

import "fmt"

type URLConflict struct {
	Str string
}

func (is URLConflict) Error() string {
	return fmt.Sprintf("error:  url: %v has already been shortened", is.Str)
}
