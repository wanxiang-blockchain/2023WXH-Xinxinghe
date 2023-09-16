package util_http

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"triple_star/util/log_formatter"
	merror "triple_star/util/util_error"
)

const (
	StatusCodeFail    = 600
	StatusCodeSuccess = 200
)

func ParseParameter(arg any, w http.ResponseWriter, r *http.Request) error {
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, arg)
	if err != nil {
		logrus.WithField("error-msg", err).
			Errorf("%s: parse parameters failed", r.URL.String())
		data := ErrorMessage(merror.ServerInternalError, "parse parameters failed", "")
		w.WriteHeader(StatusCodeFail)
		_, _ = w.Write(data)
		return err
	}
	logrus.Infof("%s: get parameters, %+v", r.URL.String(), log_formatter.TruncateLogStr(fmt.Sprintf("%+v", arg)))
	return nil
}

func ServerInternalErrorResponse(w http.ResponseWriter) {
	ErrorResponse(merror.ServerInternalError, "operation stop, server internal error", "", w)
}

func ErrorResponse(code int, desc, msg string, w http.ResponseWriter) {
	bs := ErrorMessage(code, desc, msg)
	w.WriteHeader(StatusCodeFail)
	_, _ = w.Write(bs)
}
