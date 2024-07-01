package web

type LoginSMSReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
