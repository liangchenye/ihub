package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/controllers"
)

const (
	// repo is more like a ftp protocal, we have to use a different prefix
	repoMetaPrefix = "/repometa"
)

func init() {
	if err := RegisterRouter(repoMetaPrefix, RepoMetaNameSpace()); err != nil {
		logs.Error("Failed to register router: '%s'.", repoMetaPrefix)
	} else {
		logs.Debug("Register router '%s' registered.", repoMetaPrefix)
	}
}

// RepoMetaNameSpace defines the app repo metadata router
func RepoMetaNameSpace() *beego.Namespace {
	ns := beego.NewNamespace(repoMetaPrefix,
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),
		beego.NSRouter("/*", &controllers.RepoMeta{}, "get:GetPackageMeta"),
	)

	return ns
}
