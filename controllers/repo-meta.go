package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/isula/ihub/models"
)

// RepoMeta defines the repo metadata operations
// this is used for application packages
type RepoMeta struct {
	beego.Controller
}

// GetPackageMeta get the package metadata
func (r *RepoMeta) GetPackageMeta() {
	url := r.Ctx.Input.Param(":splat")

	logs.Debug("GetPackageMeta '%s'", url)
	// TODO: meta for directory is not finish

	pkg, err := models.QueryPkgByName(filepath.Dir(url), filepath.Base(url))
	if err != nil {
		CtxErrorWrap(r.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Cannot query the file '%s'.", url))
		return
	} else if pkg == nil {
		CtxErrorWrap(r.Ctx, http.StatusNotFound, err, fmt.Sprintf("Cannot find the file '%s'.", url))
		return
	}

	CtxSuccessWrap(r.Ctx, http.StatusOK, models.PkgInfoNew(pkg), nil)
}
