package httpx

type RequestHeader struct {
}

type Request struct {
	Method string

	URI    string
	Header map[string]string

	IsJson bool
}

type RequestClient struct {
	ClientHost string

	Header    map[string][]string
	URLParams map[string][]string
	//Cookies

	JSONBody *interface{}

	request *Request
}

type Response struct {
	Code string
	Body string

	JSONBody *interface{}
}
