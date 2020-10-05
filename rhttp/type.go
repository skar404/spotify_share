package rhttp

type Result struct {
	Code int
	Url  string
}

type ResultByte struct {
	Result
	Body []byte
}

type ResultRaw struct {
	Result
	Body string
}

type ResultJson struct {
	Result
	Body interface{}
	Raw  string
}
