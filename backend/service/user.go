package service

import (
	"crypto/md5"
	"encoding/hex"
	"time"
	"triple_star/dao"
	"triple_star/service/parameter"
	merror "triple_star/util/util_error"
)

type user struct{}

var User = &user{}

func (u *user) Register(para *parameter.RegisterReq) (*parameter.GeneralOperationResp, error) {
	record := dao.User.GetByAddr(para.Addr)
	if record.Addr != "" {
		return nil, &merror.Error{Code: userHasRegistered, Desc: "the user has registered"}
	}

	sum := md5.New().Sum([]byte(para.Password))
	pwd := hex.EncodeToString(sum)
	info := &dao.UserInfo{
		Name:     para.Name,
		Addr:     para.Addr,
		Password: pwd,
		Gender:   para.Gender,
		Memo:     para.Memo,
	}
	err := dao.User.Insert(info)
	if err != nil {
		return nil, &merror.Error{
			Code: dbError,
			Desc: "insert into database failed",
		}
	}
	return &parameter.GeneralOperationResp{Success: true}, nil
}

func (u *user) Login(para *parameter.LoginReq) (*parameter.LoginResp, error) {
	info := dao.User.GetByAddr(para.Addr)
	if info.Name == "" {
		return nil, &merror.Error{
			Code: userNotFound,
			Desc: "user not found",
		}
	}

	sum := md5.New().Sum([]byte(para.Password))
	pwd := hex.EncodeToString(sum)
	if info.Password != pwd {
		return nil, &merror.Error{
			Code: passwordOrAddrError,
			Desc: "user password or addr error",
		}
	}

	sess, err := Session.New(info.Name, info.Addr)
	if err != nil {
		return nil, err
	}
	resp := &parameter.LoginResp{
		Token:  sess.Token,
		Expiry: sess.Expiry.Format(time.DateTime),
	}
	return resp, nil
}

func (u *user) Logout(token string) (*parameter.GeneralOperationResp, error) {
	_ = Session.Destroy(token)
	return &parameter.GeneralOperationResp{Success: true}, nil
}
