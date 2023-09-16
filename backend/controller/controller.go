package controller

import (
	"encoding/json"
	"net/http"
	"triple_star/util/util_http"
)

func handleResp(val any, w http.ResponseWriter) {
	bs, _ := json.Marshal(val)
	w.WriteHeader(util_http.StatusCodeSuccess)
	_, _ = w.Write(bs)
}
