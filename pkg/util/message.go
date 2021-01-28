package util

import "fmt"

type Message struct {
	sep string
	msg string
}

func NewMessage(sep string) *Message {
	return &Message{sep: sep}
}

func (m *Message) Add(msg string) {
	m.msg = fmt.Sprintf("%s%s%s", m.msg, m.sep, msg)
}

func (m *Message) String() string {
	return m.msg
}
