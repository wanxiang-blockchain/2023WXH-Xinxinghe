package dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"triple_star/util/util_fn"
)

type DataLabelInfo struct {
	LabelId uint
	DataId  uint
}

var DataLabel = &DataLabelInfo{}

func (di *DataLabelInfo) TableName() string {
	return "data_label"
}
func (di *DataLabelInfo) DB() *sqlx.DB {
	return db.LoadByModel(di)
}

func (di *DataLabelInfo) createTable() {
	sqlStr := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
		    label_id int UNSIGNED NOT NULL,
		    data_id int UNSIGNED NOT NULL,
		    PRIMARY KEY(label_id, data_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		di.TableName())
	_, err := di.DB().Exec(sqlStr)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			Errorf("create table %s failed", di.TableName())
	}
}

func (di *DataLabelInfo) Insert(info *DataLabelInfo) error {
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

func (di *DataLabelInfo) GetByIdList(ids []uint, typ int) []*DataLabelInfo {
	var (
		query string
		args  []any
		err   error
	)
	switch typ {
	case 1: // data id
		query, args, err = sqlx.In("WHERE data_id IN (?)", ids)
	case 2: // label id
		query, args, err = sqlx.In("WHERE label_id IN (?)", ids)
	}
	if err != nil {
		return nil
	}

	sqlStr := fmt.Sprintf("SELECT %s FROM %s %s", fields(di), di.TableName(), di.DB().Rebind(query))
	var infos []*DataLabelInfo
	err = di.DB().Select(&infos, sqlStr, args...)
	if err != nil {
		logrus.WithField("err-msg", fmt.Sprintf("[db-error] %s", err.Error())).
			WithField("sql", sqlStr).WithField("args", util_fn.JsonString(args)).
			Errorf("query from %s failed", di.TableName())
	}
	return infos
}
