package storage

import (
	"fmt"
	"io"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/storage/driver"
)

const repoPrefix = "repo"

// HeadRepoURL return the repo url stat
// TODO we need to get user in ctx, or setting in config
func HeadRepoURL(ctx *context.Context, url string) (driver.FileInfo, error) {
	storagePath := fmt.Sprintf("%s/%s", repoPrefix, url)
	logs.Debug("Head '%s'.", storagePath)

	return Driver().Stat(*ctx, storagePath)
}

// GetRepoPackage gets the blob data
// TODO we need to get user in ctx, or setting in config
func GetRepoPackage(ctx *context.Context, url string) ([]byte, error) {
	storagePath := fmt.Sprintf("%s/%s", repoPrefix, url)
	logs.Debug("GetRepoPackage '%s'.", storagePath)

	return Driver().GetContent(*ctx, storagePath)
}

// PutPackage put the blob data to a repo by a name
func PutPackage(ctx *context.Context, reponame string, pkgname string, data []byte) error {
	storagePath := fmt.Sprintf("%s/%s/%s", repoPrefix, reponame, pkgname)
	logs.Debug("PutPackage '%s'.", storagePath)

	return Driver().PutContent(*ctx, storagePath, data)
}

// PutPackageFromReader writes the package data to a repo from a reader stream
func PutPackageFromReader(ctx *context.Context, reponame string, pkgname string, r io.Reader) (int64, error) {
	storagePath := fmt.Sprintf("%s/%s/%s", repoPrefix, reponame, pkgname)
	logs.Debug("PutPackageFromReader '%s'.", storagePath)

	w, err := Driver().Writer(*ctx, storagePath, false)
	if err != nil {
		return 0, err
	}
	defer w.Close()

	return io.Copy(w, r)
}

// ListRepoDir lists the content inside a repo directory
func ListRepoDir(ctx *context.Context, url string) ([]string, error) {
	storagePath := fmt.Sprintf("%s/%s", repoPrefix, url)
	logs.Debug("List '%s'.", storagePath)

	raw, err := Driver().List(*ctx, storagePath)
	if err != nil {
		return nil, err
	}

	return raw, nil
	//	var clean []string
	//	for _, v := range raw {
	//		logs.Debug("Get raw ", v)
	//		clean = append(clean, filepath.Base(v)+"/")
	//	}
	//	return clean, nil
}
