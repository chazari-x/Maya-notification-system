package model

type MsgStruct struct {
	MsgType string `json:"msg_type,omitempty"`
	MsgTo   string `json:"msg_to"`
	Msg     string `json:"msg"`
}
