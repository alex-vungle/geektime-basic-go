package simpleim

type Event struct {
	Msg Message
	// 接收者
	Receiver int64
}

const eventName = "simple_im_msg"
