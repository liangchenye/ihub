package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/controllers"
)

const (
	repoPrefix = "/repo"
)

func init() {
	if err := RegisterRouter(repoPrefix, RepoNameSpace()); err != nil {
		logs.Error("Failed to register router: '%s'.", repoPrefix)
	} else {
		logs.Debug("Register router '%s' registered.", repoPrefix)
	}
}

// RepoNameSpace defines the app repo router
func RepoNameSpace() *beego.Namespace {
	ns := beego.NewNamespace(repoPrefix,
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),
		beego.NSRouter("/*", &controllers.Repo{}, "get:GetPackageOrDir;put:PutPackage"),
	)

	return ns
}
