package util_http

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type Error struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"`
}

func ErrorMessage(code int, desc string, msg string) []byte {
	data := &Error{
		Code: code,
		Desc: desc,
		Msg:  msg,
	}
	bs, err := json.Marshal(data)
	if err != nil {
		logrus.WithField("error-msg", err).
			Errorf("marshal error data failed, code %d", code)
	}
	return bs
}
