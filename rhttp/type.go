package rhttp

type Response struct {
	Code int
	Url  string
}

type ResponseByte struct {
	Response
	Body []byte
}

type ResponseRaw struct {
	Response
	Body string
}

type ResponseJson struct {
	Response
	Body interface{}
	Raw  string
}
