package util

import (
	"bytes"
)

type Message struct {
	sep string
	msg bytes.Buffer
}

func NewMessage(sep string) *Message {
	return &Message{sep: sep}
}

func (m *Message) Add(msg string) {
	m.msg.WriteString(m.sep)
	m.msg.WriteString(msg)
}

func (m *Message) String() string {
	return m.msg.String()
}
