package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"time"
	"triple_star/service"
	"triple_star/util/log_formatter"
	"triple_star/util/util_http"
)

func StopPanic(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer log_formatter.BlankPanic()
		h.ServeHTTP(w, r)
	}
}

func HandleException(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				util_http.ServerInternalErrorResponse(w)
				time.Sleep(3 * time.Second)
				logrus.WithFields(logrus.Fields{
					"panic-msg":   err,
					"stack trace": string(debug.Stack()),
				}).Panicf("http panic")
			}
		}()
		h.ServeHTTP(w, r)
	}
}

func HandleCorsRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("URL:%s Method:%s.", r.URL.Path, r.Method)
		if r.Header.Get("Origin") != "" {
			w.Header().Set("Accept-Language", "zh-CN,zh;q=0.9")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,X-Requested-With,X-Token,X-Expiry")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Expires", "0")
			w.Header().Set("Content-Type", "application/octet-stream;charset=UTF-8")
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			_, _ = w.Write(nil)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/triple_star/user/register":
		case r.URL.Path == "/v1/triple_star/user/login":
		default:
			token := r.Header.Get("X-Token")
			sess, ok := service.Session.Exist(token)
			if !ok {
				util_http.ErrorResponse(1005, "unauthorized user", "", w)
				return
			}
			w.Header().Set("X-Token", sess.Token)
			w.Header().Set("X-Expiry", sess.Expiry.Format(time.DateTime))
		}
		h.ServeHTTP(w, r)
	}
}
