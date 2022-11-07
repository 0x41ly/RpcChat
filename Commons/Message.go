package commons

import "time"

type Message struct {
	User    string
	Target  string
	Msg     string
	MsgDate time.Time
}
