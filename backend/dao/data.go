package dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"triple_star/util/util_fn"
)

type DataInfo struct {
	ID       uint `db:"id"`
	Hash     string
	Name     string
	Addr     string
	Mark     string // json: k-v map
	File     string
	Memo     string
	Category string
	CreateAt string
	Price    float64
}

var Data = &DataInfo{}

func (di *DataInfo) TableName() string {
	return "data"
}
func (di *DataInfo) DB() *sqlx.DB {
	return db.LoadByModel(di)
}

func (di *DataInfo) createTable() {
	sqlStr := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
		    id int UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
		    hash varchar(255) UNIQUE NOT NULL DEFAULT '',
		    name varchar(128) NOT NULL DEFAULT '',
		    addr varchar(64) NOT NULL DEFAULT '',
		    memo varchar(255) NOT NULL DEFAULT '',
		    mark longtext,
		    file varchar(255) NOT NULL DEFAULT '',
			category varchar(64) NOT NULL DEFAULT '',
		    create_at datetime DEFAULT CURRENT_TIMESTAMP,
		    price double NOT NULL DEFAULT 0
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		di.TableName())
	_, err := di.DB().Exec(sqlStr)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			Errorf("create table %s failed", di.TableName())
	}
}

func (di *DataInfo) Insert(info *DataInfo) error {
	cols, plh, args := insert(info)
	sqlStr := fmt.Sprintf("INSERT INTO %s %s VALUES %s", di.TableName(), cols, plh)
	_, err := di.DB().Exec(sqlStr, args...)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(args)).
			Errorf("insert into %s failed", di.TableName())
	}
	return err
}

func (di *DataInfo) GetByIdList(ids []uint) []*DataInfo {
	query, args, err := sqlx.In("WHERE id IN (?)", ids)
	if err != nil {
		return nil
	}
	sqlStr := fmt.Sprintf("SELECT id, hash,name, IFNULL(mark, '') AS mark, IFNULL(data, '') AS data  FROM %s %s", di.TableName(), di.DB().Rebind(query))
	var infos []*DataInfo
	err = di.DB().Select(&infos, sqlStr, args...)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(args)).
			Errorf("query from %s failed", di.TableName())
	}
	return infos
}

func (di *DataInfo) GetByHash(hash string) *DataInfo {
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE hash = ?", fields(di), di.TableName())
	var info DataInfo
	err := di.DB().Get(&info, sqlStr, hash)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(hash)).
			Errorf("query from %s failed", di.TableName())
	}
	return &info
}

func (di *DataInfo) GetById(id uint) *DataInfo {
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", fields(di), di.TableName())
	var info DataInfo
	err := di.DB().Get(&info, sqlStr, id)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(id)).
			Errorf("query from %s failed", di.TableName())
	}
	return &info
}

func (di *DataInfo) GetByAddr(addr string) []*DataInfo {
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE addr = ?", fields(di), di.TableName())
	var infos []*DataInfo
	err := di.DB().Get(&infos, sqlStr, addr)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(addr)).
			Errorf("query from %s failed", di.TableName())
	}
	return infos
}
