package web

type Result struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data"`
}
