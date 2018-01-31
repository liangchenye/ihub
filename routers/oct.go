package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/controllers"
)

const (
	octPrefix = "/oct"
)

func init() {
	if err := RegisterRouter(octPrefix, octNameSpace()); err != nil {
		logs.Error("Failed to register router: '%s'.", octPrefix)
	} else {
		logs.Debug("Register router '%s' registered.", octPrefix)
	}
}

// octNameSpace defines the oct router
func octNameSpace() *beego.Namespace {
	ns := beego.NewNamespace(octPrefix,
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),
		beego.NSRouter("/:repo/status", &controllers.Oct{}, "get:GetStatus"),
		beego.NSRouter("/:repo/image", &controllers.Oct{}, "get:GetImage"),
		beego.NSRouter("/:repo", &controllers.Oct{}, "post:AddOutput"),
	)

	return ns
}
