package message

type Request struct {
	Action int
	Params map[string]interface{}
}