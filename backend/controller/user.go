package controller

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"triple_star/service"
	"triple_star/service/parameter"
	merror "triple_star/util/util_error"
	"triple_star/util/util_http"
)

type user struct{}

var User = &user{}

func (u *user) Register(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseMultipartForm(32 << 1)
	logrus.WithField("data", r.Form.Encode()).Infoln("get parameter")

	gender, _ := strconv.ParseInt(r.Form.Get("gender"), 10, 32)
	arg := &parameter.RegisterReq{
		Name:     r.Form.Get("username"),
		Addr:     r.Form.Get("address"),
		Password: r.Form.Get("password"),
		Gender:   int(gender),
		Memo:     r.Form.Get("memo"),
	}
	resp, err := service.User.Register(arg)
	if err != nil {
		e := err.(*merror.Error)
		logrus.WithField("code", e.Code).WithField("err-msg", err).
			Infoln("register failed")
		util_http.ErrorResponse(e.Code, e.Desc, e.Message, w)
		return
	}
	handleResp(resp, w)
}

func (u *user) Login(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseMultipartForm(32 << 1)
	logrus.WithField("data", r.Form.Encode()).Infoln("get parameter")

	arg := &parameter.LoginReq{
		Addr:     r.Form.Get("address"),
		Password: r.Form.Get("password"),
	}
	resp, err := service.User.Login(arg)
	if err != nil {
		e := err.(*merror.Error)
		logrus.WithField("code", e.Code).WithField("err-msg", err).
			Infoln("register failed")
		util_http.ErrorResponse(e.Code, e.Desc, e.Message, w)
		return
	}
	w.Header().Set("X-Token", resp.Token)
	w.Header().Set("X-Expiry", resp.Expiry)
	handleResp(resp, w)
}

func (u *user) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Token")
	resp, err := service.User.Logout(token)
	if err != nil {
		e := err.(*merror.Error)
		logrus.WithField("code", e.Code).WithField("err-msg", err).
			Infoln("logout failed")
		util_http.ErrorResponse(e.Code, e.Desc, e.Message, w)
		return
	}
	handleResp(resp, w)
}
