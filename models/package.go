package models

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// PkgRepo defines the meta data of a package repository.
// For example, the create data, the star counts and the etc..
type PkgRepo struct {
	ID int64 `orm:"column(id);auto"`
	// More like a directory here
	Name        string `orm:"unique;column(name);size(255);null"`
	Star        int64  `orm:"column(star);null"`
	Description string `orm:"column(description);null"`
}

// Pkg defines the pkg pkg struct.
type Pkg struct {
	ID int64 `orm:"column(id);auto"`
	// Name is not unique
	Name      string `orm:"column(name);size(255);null"`
	Size      int64  `orm:"column(size);null"`
	RepoID    int64  `orm:"column(repo_id);null"`
	Downloads int64  `orm:"column(downloads);null"`

	// rpm, iso, ...
	Type string `orm:"column(type);size(15);null"`
}

var pkgModels = []interface{}{
	new(PkgRepo),
	new(Pkg),
}

func init() {
	orm.RegisterModel(pkgModels...)
}

const (
	queryPkgListByRepoName = `select * from pkg join pkg_repo pr
	     on pkg.repo_id=pr.id where pr.name=?`
	queryPkgByID   = `select * from pkg where repo_id=? and name=?`
	queryPkgByName = `select * from pkg join pkg_repo pr
	     on pkg.repo_id=pr.id where pr.name=? and pkg.name=?`
)

// QueryPkgListByRepoName returns the package list by 'reponame'
func QueryPkgListByRepoName(reponame string) ([]Pkg, error) {
	var pkgs []Pkg
	_, err := orm.NewOrm().Raw(queryPkgListByRepoName, reponame).QueryRows(&pkgs)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryPkgListByRepoName] %s", err)
		return nil, err
	}

	return pkgs, nil
}

// AddPkgRepo adds a repo to the database
func AddPkgRepo(reponame string) (*PkgRepo, error) {
	repo := &PkgRepo{}

	if err := orm.NewOrm().QueryTable("pkg_repo").
		Filter("Name__exact", reponame).One(repo); err == nil {
		logs.Debug("[AddPkgRepo] repo '%s' is exist.", reponame)
		return repo, nil
	} else if err != orm.ErrNoRows {
		logs.Error("[AddPkgRepo] fail to find repo '%s': %v", reponame, err)
		return nil, err
	}

	repo.Name = reponame
	if _, err := orm.NewOrm().Insert(repo); err != nil {
		logs.Error("[AddPkgRepo] fail to insert repo '%s': %v", reponame, err)
		return nil, err
	}

	return repo, nil
}

// QueryPkgByID returns a package by 'repoid, name'
func QueryPkgByID(repoid int64, name string) (*Pkg, error) {
	var pkgs []Pkg

	_, err := orm.NewOrm().Raw(queryPkgByID, repoid, name).QueryRows(&pkgs)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryPkgByID] %v", err)
		return nil, err
	}

	if len(pkgs) == 0 {
		logs.Debug("[QueryPkgByID] cannot find the row.")
		return nil, nil
	}

	return &pkgs[0], nil
}

// QueryPkgByName returns a package by 'repoName, name'
func QueryPkgByName(repoName string, name string) (*Pkg, error) {
	var pkgs []Pkg

	_, err := orm.NewOrm().Raw(queryPkgByName, repoName, name).QueryRows(&pkgs)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryPkgByName] %v", err)
		return nil, err
	}

	if len(pkgs) == 0 {
		logs.Debug("[QueryPkgByName] cannot find the row.")
		return nil, nil
	}

	return &pkgs[0], nil
}

// AddPkg adds a package to the database. If the target repo is not exist, it will create a repo.
//TODO: lots of rollback
func AddPkg(reponame string, name string, size int64, pkgType string) (*Pkg, error) {
	repo, err := AddPkgRepo(reponame)
	if err != nil {
		return nil, err
	}

	if pkg, err := QueryPkgByID(repo.ID, name); err != nil {
		logs.Error("[AddPkg] %v", err)
		return nil, err
	} else if pkg != nil {
		// Already exist, TODO: update info?
		logs.Debug("[AddPkg] pkg is already exist")
		pkg.Size = size
		pkg.Type = pkgType
		if _, err := orm.NewOrm().Update(pkg); err != nil {
			return nil, err
		}
		return pkg, nil
	}

	pkg := &Pkg{}
	pkg.RepoID = repo.ID
	pkg.Name = name
	pkg.Size = size
	pkg.Type = pkgType
	if _, err := orm.NewOrm().Insert(pkg); err != nil {
		logs.Error("[AddPkg] fail to insert pkg '%s:%s': %v", reponame, name, err)
		return nil, err
	}
	return pkg, nil
}

// PkgDownloadInc increases the download count by one
// TODO: lock/unlock the db item
// TODO: add a detailed pkg download info, like ip, date?
func PkgDownloadInc(reponame string, name string) (*Pkg, error) {
	logs.Debug("[PkgDownloadInc] '%s:%s'", reponame, name)
	repo, err := AddPkgRepo(reponame)
	if err != nil {
		return nil, err
	}

	pkg, err := QueryPkgByID(repo.ID, name)
	if err != nil {
		logs.Error("[PkgDownloadInc] %v", err)
		return nil, err
	} else if pkg == nil {
		logs.Error("[PkgDownloadInc] Cannot find the target repo %s:%s", reponame, name)
		return nil, fmt.Errorf("Cannot find the target repo %s:%s", reponame, name)
	}

	// Already exist, TODO: update info?
	logs.Debug("[PkgDownloadInc] pkg is already exist, increase the downloads")
	pkg.Downloads = pkg.Downloads + 1
	if _, err := orm.NewOrm().Update(pkg); err != nil {
		return nil, err
	}
	return pkg, nil
}
