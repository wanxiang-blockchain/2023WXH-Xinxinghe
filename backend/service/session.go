package service

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/sirupsen/logrus"
	"time"
	"triple_star/dao"
	merror "triple_star/util/util_error"
)

var sessMgr *scs.SessionManager

func sessInit() {
	sessMgr = scs.New()
	sessMgr.Store = dao.GetRedisStore()
}

type session struct{}

var Session = &session{}

type SessionInfo struct {
	Token  string
	Expiry time.Time
}

func (s *session) New(username string, addr string) (*SessionInfo, error) {
	ctx, err := sessMgr.Load(context.Background(), "")
	if err != nil {
		return nil, &merror.Error{
			Code: getSessionFailed,
			Desc: merror.Wrap(err, "create session failed").Error(),
		}
	}

	sessMgr.Put(ctx, "username", username)
	sessMgr.Put(ctx, "addr", addr)
	token, expiry, _ := sessMgr.Commit(ctx)
	return &SessionInfo{Token: token, Expiry: expiry.In(time.Now().Location())}, nil
}

func (s *session) Get(token string, key string) any {
	ctx, err := sessMgr.Load(context.Background(), token)
	if err != nil {
		logrus.WithField("err-msg", err).Errorln("get session failed")
		return nil
	}
	return sessMgr.Get(ctx, key)
}

func (s *session) Exist(token string) (*SessionInfo, bool) {
	ctx, err := sessMgr.Load(context.Background(), token)
	if err != nil {
		logrus.WithField("err-msg", err).Errorln("get session failed")
		return nil, false
	}

	b, found, err := sessMgr.Store.Find(token)
	if err != nil {
		logrus.WithField("err-msg", err).Errorln("find Token failed")
		return nil, false
	}
	if found {
		dt := time.Now().UTC()
		expiry, _, _ := sessMgr.Codec.Decode(b)
		if expiry.Sub(dt) < 30*time.Minute {
			_ = sessMgr.RenewToken(ctx)
			token, expiry, _ = sessMgr.Commit(ctx)
			return &SessionInfo{Token: token, Expiry: expiry.In(time.Now().Location())}, true
		}
		return &SessionInfo{Token: token, Expiry: expiry.In(time.Now().Location())}, true
	}
	return nil, false
}

func (s *session) Destroy(token string) error {
	ctx, err := sessMgr.Load(context.Background(), token)
	if err != nil {
		return &merror.Error{
			Code: getSessionFailed,
			Desc: merror.Wrap(err, "create session failed").Error(),
		}
	}

	err = sessMgr.Destroy(ctx)
	if err != nil {
		return &merror.Error{
			Code: destroySessionFailed,
			Desc: merror.Wrap(err, "destroy session failed").Error(),
		}
	}
	return nil
}
