package dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"triple_star/util/util_fn"
)

type UserInfo struct {
	ID       uint `db:"id"`
	Name     string
	Addr     string
	Password string
	Gender   int
	Memo     string
}

var User = &UserInfo{}

func (u *UserInfo) TableName() string {
	return "user"
}
func (u *UserInfo) DB() *sqlx.DB {
	return db.LoadByModel(u)
}

func (u *UserInfo) createTable() {
	sqlStr := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
		    id int UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
		    name varchar(64) NOT NULL DEFAULT '',
		    addr varchar(64) UNIQUE NOT NULL DEFAULT '',
		    password varchar(64) NOT NULL DEFAULT '',
		    gender tinyint(1) NOT NULL DEFAULT 1,
		    memo varchar(255) NOT NULL DEFAULT ''
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`, u.TableName())
	_, err := u.DB().Exec(sqlStr)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			Errorf("create table %s failed", u.TableName())
	}
}

func (u *UserInfo) Insert(info *UserInfo) error {
	cols, plh, args := insert(info)
	sqlStr := fmt.Sprintf("INSERT INTO %s %s VALUES %s", u.TableName(), cols, plh)
	_, err := u.DB().Exec(sqlStr, args...)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(args)).
			Errorf("insert into %s failed", u.TableName())
	}
	return err
}

func (u *UserInfo) GetByAddr(addr string) *UserInfo {
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE addr = ?", fields(u), u.TableName())
	var info UserInfo
	err := u.DB().Get(&info, sqlStr, addr)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(addr)).
			Errorf("query from %s failed", u.TableName())
	}
	return &info
}

func (u *UserInfo) Update(id uint, info *UserInfo) error {
	query, args := update(info)
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", u.TableName(), query)
	args = append(args, id)
	_, err := u.DB().Exec(sqlStr, args...)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(args)).
			Errorf("update %s failed", u.TableName())
	}
	return err
}

func (u *UserInfo) Delete(id uint) error {
	sqlStr := fmt.Sprintf("DELETE FROM %s where id = ?", u.TableName())
	_, err := u.DB().Exec(sqlStr, id)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(id)).
			Errorf("delete from %s failed", u.TableName())
	}
	return err
}
