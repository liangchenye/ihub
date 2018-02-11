package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/context"

	"github.com/isula/ihub/audit"
)

//TODO: more logs info

// CtxErrorWrap wraps the error http message
func CtxErrorWrap(ctx *context.Context, code int, err error, msg string) {
	ctx.Output.SetStatus(code)
	ctx.Output.Body([]byte(msg))

	audit.Error(ctx, code, err, msg)
}

// CtxSuccessWrap wraps the success http message
func CtxSuccessWrap(ctx *context.Context, code int, result interface{}, header map[string]string) {
	ctx.Output.SetStatus(code)
	for n, v := range header {
		ctx.Output.Header(n, v)
	}
	output, _ := json.Marshal(result)
	ctx.Output.Body(output)

	audit.Info(ctx, code)
}

// CtxDataWrap wraps the http data steam
func CtxDataWrap(ctx *context.Context, code int, result []byte, header map[string]string) {
	ctx.Output.SetStatus(code)
	for n, v := range header {
		ctx.Output.Header(n, v)
	}
	ctx.Output.Body(result)

	audit.Info(ctx, code)
}
