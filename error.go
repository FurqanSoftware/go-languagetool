package languagetool

import "strconv"

type Error struct {
	StatusCode int
}

func (e Error) Error() string {
	return "langaugetool: status " + strconv.Itoa(e.StatusCode)
}
