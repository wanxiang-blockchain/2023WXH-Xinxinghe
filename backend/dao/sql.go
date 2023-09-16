package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strings"
	"time"
	"triple_star/config"
)

const (
	ConnMaxLifetime = 25
	MaxOpenConns    = 10
	MaxIdleConns    = 6
)

// dbConfig options
type dbConfig struct {
	DbName       string
	DSN          string
	MaxIdleConns int           // SetMaxIdleConns default=6
	MaxOpenConns int           // SetMaxOpenConns default=10
	MaxLifetime  time.Duration // SetConnMaxLifetime default=25*60 seconds
}

// NewDBConfig default config
func NewDBConfig(name string, dsn string) *dbConfig {
	return &dbConfig{
		DbName:       name,
		DSN:          dsn,
		MaxIdleConns: MaxIdleConns,
		MaxOpenConns: MaxOpenConns,
		MaxLifetime:  ConnMaxLifetime * time.Minute,
	}
}
func (c *dbConfig) WithMaxIdleConns(n int) *dbConfig {
	c.MaxIdleConns = n
	return c
}
func (c *dbConfig) WithMaxOpenConns(n int) *dbConfig {
	c.MaxOpenConns = n
	return c
}
func (c *dbConfig) WithMaxLifetime(du time.Duration) *dbConfig {
	c.MaxLifetime = du
	return c
}

// camel case to snake case
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// sqlx.DB handle
type sdb map[string]*sqlx.DB

var db sdb

func (p sdb) Load(name string) *sqlx.DB {
	return p[name]
}
func (p sdb) Add(conf *dbConfig) {
	msdb, err := sqlx.Open("mysql", conf.DSN)
	if err != nil {
		logrus.WithField("panic-msg", err).
			Panicln("add database failed")
	}

	msdb.MapperFunc(toSnakeCase)
	msdb.SetConnMaxLifetime(conf.MaxLifetime)
	msdb.SetMaxOpenConns(conf.MaxOpenConns)
	msdb.SetMaxIdleConns(conf.MaxIdleConns)

	p[conf.DbName] = msdb
}
func (p sdb) LoadByModel(model Tabler) *sqlx.DB {
	return modelDBList[model.TableName()]
}

type Tabler interface {
	TableName() string
}

func getDSN(name string) string {
	conf := config.Config.Mysql
	prefix := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		conf.UserName, conf.Password, conf.Host, conf.Port)
	suffix := "?charset=utf8mb4&parseTime=True&loc=Local"
	return prefix + name + suffix
}

var modelDBList map[string]*sqlx.DB

func createDatabase(dbModelList map[string][]interface{}) {
	// create destination database if not exist
	nameList := make(map[string]struct{}, len(dbModelList))
	for k := range dbModelList {
		nameList[k] = struct{}{}
	}

	dsn := getDSN("")
	mdb, err := sql.Open("mysql", dsn)
	defer mdb.Close()
	if err != nil {
		logrus.WithField("error-msg", err).
			Panicf("%s open failed", dsn)
	}

	for k := range nameList {
		_, err = mdb.Exec("CREATE DATABASE IF NOT EXISTS " + k)
		if err != nil {
			logrus.WithField("error-msg", err).
				Panicln("create database failed")
		}
	}
}

func sqlInit() {
	db = make(sdb)
	modelDBList = make(map[string]*sqlx.DB)

	// the tables which need register database to.
	// the tables include all models
	var dbModelList = map[string][]interface{}{
		"triple_star": {User, Data, DataLabel, DataUserBuy, Label},
	}

	// create database
	createDatabase(dbModelList)

	// register db to models
	for name, lst := range dbModelList {
		conf := NewDBConfig(name, getDSN(name))
		db.Add(conf)
		for _, m := range lst {
			modelDBList[m.(Tabler).TableName()] = db.Load(name)
		}
	}

	createTable()
}

// create table

func createTable() {
	type table interface {
		createTable()
	}
	var tables = []table{
		User,
		Data,
		DataLabel,
		DataUserBuy,
		Label,
	}
	for _, t := range tables {
		t.createTable()
	}
}

// ------------------ helper function --------------------------

func getFieldName(field reflect.StructField) string {
	if v := field.Tag.Get("db"); v != "" {
		return v
	}
	return toSnakeCase(field.Name)
}

func fields(val any) string {
	v := reflect.ValueOf(val).Elem()
	vt := v.Type()

	var (
		ret string
		sep string
	)
	for i := 0; i < v.NumField(); i++ {
		name := getFieldName(vt.Field(i))
		ret += sep + name
		sep = ", "
	}
	return ret
}

func insert(val any) (string, string, []any) {
	v := reflect.ValueOf(val).Elem()
	vt := v.Type()

	var (
		cols = "("
		plh  = "("
		sep  string
		args = make([]any, 0)
	)

	for i := 0; i < v.NumField(); i++ {
		if v.IsZero() {
			continue
		}

		name := getFieldName(vt.Field(i))
		cols += sep + name
		plh += sep + "?"
		args = append(args, v.Field(i).Interface())
		sep = ", "
	}

	cols += ")"
	plh += ")"
	return cols, plh, args
}

func update(val any) (string, []any) {
	v := reflect.ValueOf(val).Elem()
	vt := v.Type()

	var (
		query string
		sep   string
		args  = make([]any, 0)
	)

	for i := 0; i < v.NumField(); i++ {
		if v.IsZero() {
			continue
		}

		name := getFieldName(vt.Field(i))
		query += sep + fmt.Sprintf("%s = ?", name)
		args = append(args, v.Field(i).Interface())
		sep = ", "
	}
	return query, args
}
