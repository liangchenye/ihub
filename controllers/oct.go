package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/isula/ihub/storage"
)

// Oct defines the oct operations
type Oct struct {
	beego.Controller
}

// GetStatus gets the detailed status
func (o *Oct) GetStatus() {
	repo := o.Ctx.Input.Param(":repo")

	logs.Debug("GetStatus '%s'", repo)

	data, err := storage.GetOctStatus(o.Ctx, repo)
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find the status of '%s'.", repo))
		return
	}
	size := len(data)
	header := make(map[string]string)
	header["Content-Length"] = fmt.Sprint(size)
	CtxDataWrap(o.Ctx, http.StatusOK, data, header)
	return
}

// GetImage gets the image for better user experience
func (o *Oct) GetImage() {
	repo := o.Ctx.Input.Param(":repo")
	logs.Debug("GetImage '%s'", repo)

	data, err := storage.GetOctImage(o.Ctx, repo)
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find the image of '%s'.", repo))
		return
	}
	size := len(data)
	header := make(map[string]string)
	header["Content-Length"] = fmt.Sprint(size)
	header["Content-Type"] = "image/svg+xml"
	CtxDataWrap(o.Ctx, http.StatusOK, data, header)
	return
}

// AddOutput adds the certification outputs
func (o *Oct) AddOutput() {
	repo := o.Ctx.Input.Param(":repo")
	logs.Debug("AddOutput '%s'", repo)

	file, _, err := o.Ctx.Request.FormFile("file")
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusBadRequest, nil, fmt.Sprintf("Cannot find the upload output info '%s'.", repo))
		return
	}
	data, _ := ioutil.ReadAll(file)
	file.Close()

	err = storage.AddOctOutput(o.Ctx, repo, data)
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, nil, fmt.Sprintf("Fail to add output info '%s'.", repo))
		return
	}
	CtxSuccessWrap(o.Ctx, http.StatusOK, "ok", nil)
}
