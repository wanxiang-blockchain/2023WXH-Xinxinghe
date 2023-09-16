package config

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"reflect"
	merror "triple_star/util/util_error"
)

const (
	iniConfName = "triple-star.ini"

	// default config value
	mysqlDefaultHost = "127.0.0.1"
	mysqlDefaultPort = 3306
	mysqlDefaultUser = "root"
	mysqlDefaultPwd  = "root"

	redisDefaultHost = "127.0.0.1"
	redisDefaultPort = 6379

	httpDefaultPort = 8080

	// ini file config key
	// we can get value by the key
	iniMysqlSectionName = "mysql"
	iniHttpSectionName  = "http"
	iniRedisSectionName = "redis"
)

// config struct ----------------------------------------------

type MysqlConfig struct {
	Host     string `name:"host"`
	Port     uint64 `name:"port"`
	UserName string `name:"user"`
	Password string `name:"password"`
}
type RedisConfig struct {
	Host string `name:"host"`
	Port uint64 `name:"port"`
}
type HttpConfig struct {
	Port uint64 `name:"port"`
}

// export var
var Config = struct {
	Mysql MysqlConfig
	Redis RedisConfig
	Http  HttpConfig
}{
	Mysql: MysqlConfig{
		Host:     mysqlDefaultHost,
		Port:     mysqlDefaultPort,
		UserName: mysqlDefaultUser,
		Password: mysqlDefaultPwd,
	},
	Redis: RedisConfig{
		Host: redisDefaultHost,
		Port: redisDefaultPort,
	},
	Http: HttpConfig{
		Port: httpDefaultPort,
	},
}

// getConfig use reflect to assign value to struct
func getConfig(i interface{}, section *ini.Section) {
	val := reflect.ValueOf(i).Elem()
	valType := val.Type()
	for i := 0; i < valType.NumField(); i++ {
		// get tag
		typ := valType.Field(i).Type.Kind()
		switch typ {
		case reflect.String:
			newVal := section.Key(valType.Field(i).Tag.Get("name")).String()
			if len(newVal) != 0 {
				val.Field(i).SetString(newVal)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			newVal, err := section.Key(valType.Field(i).Tag.Get("name")).Int64()
			if err == nil {
				val.Field(i).SetInt(newVal)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			newVal, err := section.Key(valType.Field(i).Tag.Get("name")).Uint64()
			if err == nil {
				val.Field(i).SetUint(newVal)
			}
		case reflect.Bool:
			{
				newVal, err := section.Key(valType.Field(i).Tag.Get("name")).Bool()
				if err == nil {
					val.Field(i).SetBool(newVal)
				}
			}
		default:
		}
	}
}

func GetAppDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		err = merror.Wrap(err, "get dir error")
		return "", err
	}
	wd, err := os.Getwd()
	if err != nil {
		err = merror.Wrap(err, "get wd error")
		return "", err
	}
	if dir != wd {
		return wd, nil
	}
	return dir, nil
}

func setIniFileConfig() {
	dir, err := GetAppDir()
	if err != nil {
		logrus.WithField("error-msg", err).
			Errorln("getApp dir err")
		return
	}
	conf := dir
	conf += string(os.PathSeparator)
	conf += iniConfName
	logrus.Infoln("start load config:", conf)

	iniFile, err := ini.Load(conf)
	if err != nil {
		logrus.WithField("error-msg", err).
			Errorln(" fail to load config file")
		return
	}
	// set config
	getConfig(&Config.Mysql, iniFile.Section(iniMysqlSectionName))
	getConfig(&Config.Http, iniFile.Section(iniHttpSectionName))
	getConfig(&Config.Redis, iniFile.Section(iniRedisSectionName))
}

func Init() {
	setIniFileConfig()
}
