package domain

import "time"

type Msg struct {
	Id      int64
	Content string
	Ctime   time.Time
}
