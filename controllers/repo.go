package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/isula/ihub/models"
	"github.com/isula/ihub/storage"
	"github.com/isula/ihub/storage/driver"
)

// Repo defines the repo operations
// this is used for application packages
type Repo struct {
	beego.Controller
}

// GetPackageOrDir get the package list or a package
func (r *Repo) GetPackageOrDir() {
	url := r.Ctx.Input.Param(":splat")

	logs.Debug("GetPackageOrDir '%s'", url)
	info, err := storage.HeadRepoURL(r.Ctx, url)
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); ok {
			CtxErrorWrap(r.Ctx, http.StatusNotFound, err, fmt.Sprintf("Cannot find the file '%s'.", url))
		} else {
			CtxErrorWrap(r.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Cannot query the file '%s'.", url))
		}
		return
	}

	if info.IsDir() {
		logs.Debug("'%s' is a directory", url)
		files, _ := storage.ListRepoDir(r.Ctx, url)
		// return the template view
		r.Data["files"] = files
		return
	}

	reponame := filepath.Dir(url)
	pkgname := filepath.Base(url)
	// TODO In some cases, we add pkg simple by 'rsync' or 'scp'. It should have a strick rule to disable this operation.
	//  Or we need to provide a 'sync' api.
	pkg, err := models.QueryPkgByName(reponame, pkgname)
	if err != nil {
		CtxErrorWrap(r.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Cannot query the file '%s'.", url))
		return
	}

	data, err := storage.GetRepoPackage(r.Ctx, url)
	if err != nil {
		CtxErrorWrap(r.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Cannot find the file '%s'.", url))
		return
	}
	size := len(data)

	if pkg == nil {
		if _, err := models.AddPkg(reponame, pkgname, int64(size), ""); err != nil {
			CtxErrorWrap(r.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Cannot add the file '%s'.", url))
			return
		}
	}

	if _, err := models.PkgDownloadInc(filepath.Dir(url), filepath.Base(url)); err != nil {
		logs.Error("Fail to record the download information: '%v'.", err)
	}

	header := make(map[string]string)
	header["Content-Length"] = fmt.Sprint(size)
	CtxDataWrap(r.Ctx, http.StatusOK, data, header)
	return
}

// PutPackage add a package to the repo
func (r *Repo) PutPackage() {
	url := r.Ctx.Input.Param(":splat")
	logs.Debug("PutPackage '%s'", url)

	// TODO: reuse the filename, and upload to the directory
	if strings.HasSuffix(url, "/") {
		CtxErrorWrap(r.Ctx, http.StatusBadRequest, nil, fmt.Sprintf("Cannot put a directory '%s'.", url))
		return
	}

	info, err := storage.HeadRepoURL(r.Ctx, url)
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); !ok {
			CtxErrorWrap(r.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Cannot query the file '%s'.", url))
			return
		}
	}

	// add or update.

	if info != nil && info.IsDir() {
		logs.Debug("'%s' is a directory", url)
		CtxErrorWrap(r.Ctx, http.StatusBadRequest, nil, fmt.Sprintf("Cannot overwrite a directory by a file '%s'.", url))
		return
	}

	file, _, err := r.Ctx.Request.FormFile("file")
	if err != nil {
		CtxErrorWrap(r.Ctx, http.StatusBadRequest, nil, fmt.Sprintf("Cannot find the upload file '%s'.", url))
		return
	}
	defer file.Close()

	reponame := filepath.Dir(url)
	pkgname := filepath.Base(url)
	size, err := storage.PutPackageFromReader(r.Ctx, reponame, pkgname, file)
	if err != nil {
		CtxErrorWrap(r.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put package of '%s'.", url))
		return
	}

	// TODO: Get Size from reader
	if _, err := models.AddPkg(reponame, pkgname, size, ""); err != nil {
		CtxErrorWrap(r.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put package of '%s'.", url))
		// TODO: rollback
		return
	}
	CtxSuccessWrap(r.Ctx, http.StatusOK, "{}", nil)
}
