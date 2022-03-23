package internal

type Event interface {
	Prepare() []byte
}

type StringEvent struct {
	Id    string
	Event string
	Data  string
}

func (e StringEvent) Prepare() []byte {
	//var data bytes.Buffer
	//
	//if len(e.Id) > 0 {
	//
	//}

	return nil
}
