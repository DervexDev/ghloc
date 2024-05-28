package rest

type Unauthorized struct {
	Msg string
}

func (e Unauthorized) Error() string {
	return e.Msg
}
