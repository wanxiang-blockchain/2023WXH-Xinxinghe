package dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"triple_star/util/util_fn"
)

type DataUserBuyInfo struct {
	ID     uint `db:"id"`
	UserId uint
	Addr   string
	DataId uint
	BuyAt  string
}

var DataUserBuy = &DataUserBuyInfo{}

func (di *DataUserBuyInfo) TableName() string {
	return "data_user_buy"
}
func (di *DataUserBuyInfo) DB() *sqlx.DB {
	return db.LoadByModel(di)
}

func (di *DataUserBuyInfo) createTable() {
	sqlStr := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
		    id int UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		    user_id int UNSIGNED NOT NULL DEFAULT 0,
		    addr varchar(64) NOT NULL DEFAULT '',
		    data_id int UNSIGNED NOT NULL DEFAULT 0,
		    buy_at datetime DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		di.TableName())
	_, err := di.DB().Exec(sqlStr)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			Errorf("create table %s failed", di.TableName())
	}
}

func (di *DataUserBuyInfo) Insert(info *DataUserBuyInfo) error {
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

func (di *DataUserBuyInfo) GetByAddr(addr string) []*DataUserBuyInfo {
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE addr = ?", fields(di), di.TableName())
	var infos []*DataUserBuyInfo
	err := di.DB().Get(&infos, sqlStr, addr)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(addr)).
			Errorf("query from %s failed", di.TableName())
	}
	return infos
}
