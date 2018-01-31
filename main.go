package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"

	"github.com/isula/ihub/config"
	"github.com/isula/ihub/logger"
	"github.com/isula/ihub/models"
	"github.com/isula/ihub/routers"
	"github.com/isula/ihub/session"
	_ "github.com/isula/ihub/session/memory"
	"github.com/isula/ihub/storage"
	_ "github.com/isula/ihub/storage/driver/filesystem"
)

func main() {
	cfg, err := config.InitConfigFromFile("conf/isula.yml")
	if err != nil {
		return
	}
	if err := logger.InitLogger(cfg.Log); err != nil {
		logs.Warning(err)
	}

	if err := config.InitHTTP(cfg.Server); err != nil {
		logs.Critical("Error in init http: ", err)
		return
	}

	conn, err := cfg.DB.GetConnection()
	if err == nil {
		if err := models.InitDB(conn, cfg.DB.Driver, "default"); err != nil {
			logs.Critical("Error in init db: ", err)
			return
		}
	} else {
		// Don't need to have a database sometimes
		logs.Warning(err)
	}

	if err := storage.InitStorage(cfg.Storage); err != nil {
		logs.Critical("Error in init storage: ", err)
		return
	}

	if err := session.InitSession(cfg.Session); err != nil {
		logs.Critical("Error in init session: ", err)
		return
	}

	nss := router.GetNamespaces()
	for name, ns := range nss {
		logs.Debug("Namespace '%s' is enabled", name)
		beego.AddNamespace(ns)
	}

	beego.Run()
}
