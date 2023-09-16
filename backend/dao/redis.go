package dao

import (
	"fmt"
	"github.com/alexedwards/scs/redisstore"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"time"
	"triple_star/config"
)

var pool *redis.Pool

func redisInit() {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			cnf := config.Config.Redis
			addr := fmt.Sprintf("%s:%d", cnf.Host, cnf.Port)
			conn, err := redis.Dial("tcp", addr, redis.DialDatabase(1))
			if err != nil {
				logrus.WithField("err-msg", err).Panicln("init redis failed")
			}
			return conn, nil
		},
		MaxIdle:     5,
		MaxActive:   3,
		IdleTimeout: 15 * time.Minute,
	}
}

func GetRedisStore() *redisstore.RedisStore {
	if pool == nil {
		redisInit()
	}
	return redisstore.New(pool)
}
