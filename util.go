package paperweight

type Error struct {
	error
}

func (e Error) Error() string {
	//TODO implement me
	return e.error.Error()
}

func (e Error) Handle() {
	panic(e.error.Error())
}
