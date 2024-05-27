package rest

import "fmt"

var NotAuthorized = fmt.Errorf("unauthorized")

type Unauthorized struct {}

func (e Unauthorized) Error() error {
	return NotAuthorized
}
