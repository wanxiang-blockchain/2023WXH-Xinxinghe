package controller

import (
	"github.com/issue9/upload"
	"github.com/sirupsen/logrus"
	"net/http"
	"triple_star/service"
	"triple_star/service/parameter"
	merror "triple_star/util/util_error"
	"triple_star/util/util_http"
)

type data struct{}

var Data = &data{}

func (d *data) Upload(w http.ResponseWriter, r *http.Request) {
	dir := "uploads/"
	subDirFmt := "2006-01-02/"
	u, _ := upload.New(dir, subDirFmt, 31<<1, ".zip", ".txt", ".png")
	if r.Method == "POST" {
		files, _ := u.Do("files", r) // 执行上传操作
		handleResp(&parameter.UploadResp{Files: files}, w)
	}
}

func (d *data) Add(w http.ResponseWriter, r *http.Request) {
	para := &parameter.DataAddReq{}
	if err := util_http.ParseParameter(&para, w, r); err != nil {
		return
	}

	resp, err := service.Data.Add(para)
	if err != nil {
		e := err.(*merror.Error)
		logrus.WithField("code", e.Code).WithField("err-msg", err).
			Infoln("add data failed")
		util_http.ErrorResponse(e.Code, e.Desc, e.Message, w)
		return
	}
	handleResp(resp, w)
}

func (d *data) Query(w http.ResponseWriter, r *http.Request) {
	para := &parameter.DataQueryReq{}
	if err := util_http.ParseParameter(&para, w, r); err != nil {
		return
	}

	resp, err := service.Data.Query(para)
	if err != nil {
		e := err.(*merror.Error)
		logrus.WithField("code", e.Code).WithField("err-msg", err).
			Infoln("query data failed")
		util_http.ErrorResponse(e.Code, e.Desc, e.Message, w)
		return
	}
	handleResp(resp, w)
}

func (d *data) Buy(w http.ResponseWriter, r *http.Request) {
	para := struct {
		Id uint
	}{}
	if err := util_http.ParseParameter(&para, w, r); err != nil {
		return
	}

	token := r.Header.Get("X-Token")
	addr := service.Session.Get(token, "addr").(string)

	buyPara := &parameter.DataBuyReq{
		Id:   para.Id,
		Addr: addr,
	}
	resp, err := service.Data.Buy(buyPara)
	if err != nil {
		e := err.(*merror.Error)
		logrus.WithField("code", e.Code).WithField("err-msg", err).
			Infoln("buy data failed")
		util_http.ErrorResponse(e.Code, e.Desc, e.Message, w)
		return
	}
	handleResp(resp, w)
}

func (d *data) QueryOneSelf(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Token")
	addr := service.Session.Get(token, "addr").(string)
	resp, err := service.Data.QueryOneSelf(&parameter.UserDataQueryReq{Addr: addr})
	if err != nil {
		e := err.(*merror.Error)
		logrus.WithField("code", e.Code).WithField("err-msg", err).
			Infoln("query oneself data failed")
		util_http.ErrorResponse(e.Code, e.Desc, e.Message, w)
		return
	}
	handleResp(resp, w)
}
