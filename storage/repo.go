package storage

import (
	"fmt"
	"io"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/storage/driver"
)

// HeadRepoURL return the repo url stat
// TODO we need to get user in ctx, or setting in config
func HeadRepoURL(ctx *context.Context, url string) (driver.FileInfo, error) {
	logs.Debug("Head '%s'.", url)

	return Driver().Stat(*ctx, url)
}

// GetRepoPackage gets the blob data
// TODO we need to get user in ctx, or setting in config
func GetRepoPackage(ctx *context.Context, url string) ([]byte, error) {
	logs.Debug("GetRepoPackage '%s'.", url)

	return Driver().GetContent(*ctx, url)
}

// PutPackage put the blob data to a repo by a name
func PutPackage(ctx *context.Context, reponame string, pkgname string, data []byte) error {
	url := fmt.Sprintf("%s/%s", reponame, pkgname)
	logs.Debug("PutPackage '%s'.", url)

	return Driver().PutContent(*ctx, url, data)
}

// PutPackageFromReader writes the package data to a repo from a reader stream
func PutPackageFromReader(ctx *context.Context, reponame string, pkgname string, r io.Reader) (int64, error) {
	url := fmt.Sprintf("%s/%s", reponame, pkgname)
	logs.Debug("PutPackageFromReader '%s'.", url)

	w, err := Driver().Writer(*ctx, url, false)
	if err != nil {
		return 0, err
	}
	defer w.Close()

	return io.Copy(w, r)
}

// ListRepoDir lists the content inside a repo directory
func ListRepoDir(ctx *context.Context, url string) ([]string, error) {
	logs.Debug("List '%s'.", url)

	raw, err := Driver().List(*ctx, url)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
