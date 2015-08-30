package main
import (
	"time"
)

type Message struct {
	Id uint64
	Type string
	Param string
	Param2 string
	Token string
//	reqSeq int64
//	respSeq int64
	Reading Reading  // (think NoteChange) Will include a ref to it's device
	CreatedAt time.Time
}

//func NewMessage() (*Message, error) {
//	msg := new(Message)
//	msg.SetReqSeq(0)
//	return msg, nil
//}
//
//// We use a pointer when we want the change our object
//func (m *Message) SetReqSeq(seq int64) {
//	m.reqSeq = seq
//}
//
//func (m Message) GetReqSeq() int64 {
//	return m.reqSeq
//}
//
//func (m *Message) IncReqSeq() int64 {
//	m.reqSeq += 1
//	return m.reqSeq // just for convenience
//}
//
//func (m *Message) Respond() {
//	m.respSeq = m.reqSeq
//}
//
//func (m Message) InSequence() {
//	m.respSeq == m.reqSeq
//}
