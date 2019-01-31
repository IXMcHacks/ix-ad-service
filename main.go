package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/eftakhairul/ix-ad-service/handlers"
	"github.com/eftakhairul/ix-ad-service/lib"
	"github.com/sirupsen/logrus"
)

func main() {
	var logger = logrus.New()

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logger.SetFormatter(customFormatter)

	var currentDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal("something went wrong. Error: ", err)
		return
	}

	var conf *lib.Config
	conf, err = lib.LoadConfig(currentDir, "config")
	if err != nil {
		logger.Fatal("unable to load the config. Error: ", err)
		return
	}

	var handlerObj = &handlers.HandlerObj{
		Logger: logger,
		Config: conf,
	}

	logger.Info("Attempting to start HTTP Server on port: ", conf.Port)
	http.HandleFunc("/heartbeat", handlerObj.Heartbeat)
	http.HandleFunc("/ixrtb", handlerObj.Adserving)

	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)

	if err != nil {
		logger.Error("Server failed starting. Error: ", err)
	}
}
