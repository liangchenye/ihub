package audit

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

// Error records/prints the detailed error messages
func Error(ctx *context.Context, code int, err error, msg string) {
	if err != nil {
		logs.Trace("Failed to [%s] [%s] [%d]: [%v]", ctx.Input.Method(), ctx.Input.URI(), code, err)
	} else {
		logs.Trace("Failed to [%s] [%s] [%d]", ctx.Input.Method(), ctx.Input.URI(), code)
	}
}

// Info records/prints the defailed messages
func Info(ctx *context.Context, code int) {
	logs.Trace("Succeed in [%s] [%s].", ctx.Input.Method(), ctx.Input.URI())
}
