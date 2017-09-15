package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// ContainerRepo defines the container repo struct,
// will support dockerv2, ociv1.
type ContainerRepo struct {
	ID          int    `orm:"column(id);auto"`
	Name        string `orm:"unique;column(name);size(255);null"`
	Star        int    `orm:"column(star);null"`
	DownloadNum int    `orm:"column(download_num);null"`
	Description string `orm:"column(description);null"`
}

// ContainerImage defines the container image struct.
type ContainerImage struct {
	ID     int    `orm:"column(id);auto"`
	Tag    string `orm:"column(tag);size(255);null"`
	Size   int64  `orm:"column(size);null"`
	RepoID int    `orm:"column(repo_id);null"`

	// docker, oci, rkt...
	Proto        string `orm:"column(proto);size(15);null"`
	ProtoVersion string `orm:"column(proto_version);size(15);null"`
}

var containerModels = []interface{}{
	new(ContainerRepo),
	new(ContainerImage),
}

func init() {
	orm.RegisterModel(containerModels...)
}

const (
	queryContainerTagsList = `select ci.Tag from container_image ci join container_repo cr 
	     on ci.repo_id=cr.id where cr.name=? and ci.proto=? and ci.proto_version=?`
	queryContainerReposList = `select name from container_repo order by id asc`
	queryContainerImage     = `select * from container_image 
	     where repo_id=? and tag=? and proto=? and proto_version=? limit 1`
)

// QueryTagsList returns the tags list by 'reponame, proto and proto version'
func QueryTagsList(reponame string, proto string, protoVerion string) ([]string, error) {
	var tags []string
	_, err := orm.NewOrm().Raw(queryContainerTagsList, reponame, proto, protoVerion).QueryRows(&tags)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryTagsList] %s", err)
		return nil, err
	}

	return tags, nil
}

// QueryReposList returns the repos list
func QueryReposList() ([]string, error) {
	var names []string
	_, err := orm.NewOrm().Raw(queryContainerReposList).QueryRows(&names)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryNamesList] %s", err)
		return nil, err
	}

	return names, nil
}

// AddRepo adds a repo to the database
func AddRepo(reponame string) (*ContainerRepo, error) {
	repo := &ContainerRepo{}

	if err := orm.NewOrm().QueryTable("container_repo").
		Filter("Name__exact", reponame).One(repo); err == nil {
		logs.Debug("[AddRepo] repo '%s' is exist.", reponame)
		return repo, nil
	} else if err != orm.ErrNoRows {
		logs.Error("[AddRepo] fail to find repo '%s': %v", reponame, err)
		return nil, err
	}

	repo.Name = reponame
	if _, err := orm.NewOrm().Insert(repo); err != nil {
		logs.Error("[AddRepo] fail to insert repo '%s': %v", reponame, err)
		return nil, err
	}

	return repo, nil
}

// QueryImage returns a container image by 'repoid, tag, proto and proto version'
func QueryImage(repoid int, tag string, proto string, protoVerion string) (*ContainerImage, error) {
	var images []ContainerImage

	_, err := orm.NewOrm().Raw(queryContainerImage, repoid, tag, proto, protoVerion).QueryRows(&images)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryImage] %v", err)
		return nil, err
	}

	if len(images) == 0 {
		logs.Debug("[QueryImage] cannot find the row.")
		return nil, nil
	}

	return &images[0], nil
}

// AddImage adds an image to the database. If the target repo is not exist, it will create a repo.
//TODO: lots of rollback
func AddImage(reponame string, tags string, proto string, protoVerion string) (*ContainerImage, error) {
	repo, err := AddRepo(reponame)
	if err != nil {
		return nil, err
	}

	if img, err := QueryImage(repo.ID, tags, proto, protoVerion); err != nil {
		logs.Error("[AddImage] %v", err)
		return nil, err
	} else if img != nil {
		// Already exist, TODO: update info?
		logs.Debug("[AddImage] image is already exist")
		return img, nil
	}

	image := &ContainerImage{}
	image.RepoID = repo.ID
	image.Tag = tags
	image.Proto = proto
	image.ProtoVersion = protoVerion
	if _, err := orm.NewOrm().Insert(image); err != nil {
		logs.Error("[AddImage] fail to insert image '%s:%s': %v", reponame, tags, err)
		return nil, err
	}
	return image, nil
}
