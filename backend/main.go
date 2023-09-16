package main

import (
	"triple_star/config"
	"triple_star/server"
	"fmt"
	"flag"
	"github.com/chai2010/winsvc"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime/debug"
)

var (
	appPath              string
	flagServiceName      = flag.String("service-name", "triple_star", "Set service name")
	flagServiceDesc      = flag.String("service-desc", "triple_star", "Set service description")
	flagServiceInstall   = flag.Bool("service-install", false, "Install service")
	flagServiceUninstall = flag.Bool("service-remove", false, "Remove service")
	flagServiceStart     = flag.Bool("service-start", false, "Start service")
	flagServiceStop      = flag.Bool("service-stop", false, "Stop service")
	flagServiceVersion   = flag.Bool("version", false, "Print service version")
	flagServiceLog       = flag.String("log", "info", "log level: panic, fatal, error, warn, info, debug, trace")
)

//Init Server
func initWinServer() error {
	var err error
	if appPath, err = winsvc.GetAppPath(); err != nil {
		logrus.WithField("error-msg", err).
			Errorln("get app path error")
		return err
	}
	if err = os.Chdir(filepath.Dir(appPath)); err != nil {
		logrus.WithField("error-msg", err).
			Errorf("change dir error, appPath:%s", appPath)
		return err
	}
	return err
}

func StartServer() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("error-msg", r).
				WithField("stack-trace", string(debug.Stack())).
				Errorln("main stop, panic")
		}
	}()
	server.Start(*flagServiceLog)
}

func StopServer() {
	logrus.Errorln("stop server")
}

func main() {
	if initWinServer() != nil {
		logrus.Errorln("initWinServer Error")
		return
	}
	flag.Parse()

	if *flagServiceVersion {
		fmt.Printf("%s: version %s", config.AppName, config.Version)
		return
	}

	if *flagServiceInstall {
		if err := winsvc.InstallService(appPath, *flagServiceName, *flagServiceDesc); err != nil {
			logrus.WithField("error-msg", err).
				Errorf("install Service(%s,%s) error", *flagServiceName, *flagServiceDesc)
		}
		logrus.Infoln("install Service Done")
		return
	}

	if *flagServiceUninstall {
		if err := winsvc.RemoveService(*flagServiceName); err != nil {
			logrus.WithField("error-msg", err).
				Errorln("remove service error")
		}
		logrus.Infoln("remove service done")
		return
	}

	if *flagServiceStart {
		if err := winsvc.StartService(*flagServiceName); err != nil {
			logrus.WithField("error-msg", err).
				Errorln("start service error")
		}
		logrus.Infoln("start service done")
		return
	}
	if *flagServiceStop {
		if err := winsvc.StopService(*flagServiceName); err != nil {
			logrus.WithField("error-msg", err).
				Errorln("stop service error")
		}
		logrus.Infoln("stop service done")
		return
	}

	if !winsvc.InServiceMode() {
		logrus.Infoln("Run In Service")
		if err := winsvc.RunAsService(*flagServiceName, StartServer, StopServer, false); err != nil {
			logrus.WithField("error-msg", err).
				Errorln("svc.Run error")
		}
		return
	}
	StartServer()
}