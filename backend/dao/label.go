package dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"triple_star/util/util_fn"
)

type LabelInfo struct {
	ID   uint `db:"id"`
	Name string
	Memo string
}

var Label = &LabelInfo{}

func (li *LabelInfo) TableName() string {
	return "label"
}
func (li *LabelInfo) DB() *sqlx.DB {
	return db.LoadByModel(li)
}

func (li *LabelInfo) createTable() {
	sqlStr := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
		    id int UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
		    name varchar(64) NOT NULL DEFAULT '',
		    memo varchar(255) NOT NULL DEFAULT '',
		    UNIQUE(name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		li.TableName())
	_, err := li.DB().Exec(sqlStr)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			Errorf("create table %s failed", li.TableName())
	}
}

func (li *LabelInfo) Insert(info *LabelInfo) error {
	cols, plh, args := insert(info)
	sqlStr := fmt.Sprintf("INSERT INTO %s %s VALUES %s", li.TableName(), cols, plh)
	_, err := li.DB().Exec(sqlStr, args...)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(args)).
			Errorf("insert into %s failed", li.TableName())
	}
	return err
}

func (li *LabelInfo) GetByIdList(ids []uint) []*LabelInfo {
	query, args, err := sqlx.In("WHERE id IN (?)", ids)
	if err != nil {
		return nil
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s %s", fields(li), li.TableName(), li.DB().Rebind(query))
	var infos []*LabelInfo
	err = li.DB().Select(&infos, sqlStr, args...)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(args)).
			Errorf("query from %s failed", li.TableName())
	}
	return infos
}

func (li *LabelInfo) GetByNameList(names []string) []*LabelInfo {
	query, args, err := sqlx.In("WHERE name IN (?)", names)
	if err != nil {
		return nil
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s %s", fields(li), li.TableName(), li.DB().Rebind(query))
	var infos []*LabelInfo
	err = li.DB().Select(&infos, sqlStr, args...)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(args)).
			Errorf("query from %s failed", li.TableName())
	}
	return infos
}
