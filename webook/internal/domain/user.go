package domain

import "time"

type User struct {
	Id       int64  `json:"Id"`
	Email    string `json:"Email"`
	Password string `json:"-"`

	// Nickname
	Nickname string `json:"Nickname"`
	// Birthday
	Birthday string `json:"Birthday"`
	// Phone
	Phone string `json:"Phone"`
	// Bio
	Bio string `json:"AboutMe"`

	// UTC 0 的时区
	Ctime time.Time

	//Addr Address
}

//type Address struct {
//	Province string
//	Region   string
//}

//func (u User) ValidateEmail() bool {
// 在这里用正则表达式校验
//return u.Email
//}
