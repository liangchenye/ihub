package controllers

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/docker/distribution/registry/storage/driver"

	"github.com/isula/ihub/models"
	"github.com/isula/ihub/session"
	"github.com/isula/ihub/storage"
	"github.com/isula/ihub/utils"
)

// DockerV2Ping defines the /v2 status
type DockerV2Ping struct {
	beego.Controller
}

// Ping returns 'ok'
func (d *DockerV2Ping) Ping() {
	head := make(map[string]string)
	head["Content-Type"] = "application/json; charset=utf-8"
	head["Docker-Distribution-Api-Version"] = "registry/2.0"
	CtxSuccessWrap(d.Ctx, http.StatusOK, "{}", head)

}

// DockerV2Tag defines the tags operation
type DockerV2Tag struct {
	beego.Controller
}

// GetTagsList returns the tags list
func (d *DockerV2Tag) GetTagsList() {
	reponame := d.Ctx.Input.Param(":splat")
	logs.Debug("GetTagsList of '%s'.", reponame)

	repo, err := models.QueryTagsList(reponame, "docker", "v2")
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get tag list of '%s'.", reponame))
		return
	} else if len(repo) == 0 {
		CtxErrorWrap(d.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find tag list of '%s'.", reponame))
		return
	}

	CtxSuccessWrap(d.Ctx, http.StatusOK, repo, nil)
}

// DockerV2Manifest defines the manifest operations
type DockerV2Manifest struct {
	beego.Controller
}

// GetManifest returns the manifest by 'repo' and 'tag'
func (d *DockerV2Manifest) GetManifest() {
	reponame := d.Ctx.Input.Param(":splat")
	tags := d.Ctx.Input.Param(":tags")
	logs.Debug("GetManifest of '%s:%s'.", reponame, tags)

	data, err := storage.GetManifest(d.Ctx, reponame, tags, "docker", "v2")
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); ok {
			CtxErrorWrap(d.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to get manifest of '%s:%s'.", reponame, tags))
		} else {
			CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get manifest of '%s:%s'.", reponame, tags))
		}
		return
	}

	digest, _ := utils.DigestManifest(data)
	header := make(map[string]string)
	header["Content-Type"] = "application/vnd.docker.distribution.manifest.v2+json"
	header["Docker-Content-Digest"] = digest
	header["Content-Length"] = fmt.Sprint(len(data))
	CtxDataWrap(d.Ctx, http.StatusOK, data, header)
}

// PutManifest puts manifest of the 'repo:tag'
func (d *DockerV2Manifest) PutManifest() {
	reponame := d.Ctx.Input.Param(":splat")
	tags := d.Ctx.Input.Param(":tags")
	logs.Debug("PutManifest of '%s:%s'.", reponame, tags)

	data := d.Ctx.Input.CopyBody(utils.MaxSize)
	logs.Debug("The manifest is <%s>", data)
	err := storage.PutManifest(d.Ctx, reponame, tags, "docker", "v2", data)
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put manifest of '%s:%s'.", reponame, tags))
		return
	}

	//TODO: rollback the storage.. add error checks
	_, err = models.AddImage(reponame, tags, "docker", "v2")
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to add image '%s:%s' to db.", reponame, tags))
		return
	}

	digest, _ := utils.DigestManifest(data)
	header := make(map[string]string)
	header["Docker-Content-Digest"] = digest
	CtxSuccessWrap(d.Ctx, http.StatusOK, "{}", header)
}

// DeleteManifest deletes the manifest of the 'repo:tag'
func (d *DockerV2Manifest) DeleteManifest() {
	reponame := d.Ctx.Input.Param(":splat")
	tags := d.Ctx.Input.Param(":tags")
	logs.Debug("DeleteManifest of '%s:%s'.", reponame, tags)

	err := storage.DeleteManifest(d.Ctx, reponame, tags, "docker", "v2")
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete manifest of '%s:%s'.", reponame, tags))
		return
	}
	CtxSuccessWrap(d.Ctx, http.StatusOK, fmt.Sprintf("Succeed in deleting manifest of '%s:%s'.", reponame, tags), nil)
}

// DockerV2Blob defines the blob operations
type DockerV2Blob struct {
	beego.Controller
}

// HeadBlob queries the blob info
func (d *DockerV2Blob) HeadBlob() {
	reponame := d.Ctx.Input.Param(":splat")
	digest := d.Ctx.Input.Param(":digest")
	logs.Debug("HeadBlob of '%s:%s'.", reponame, digest)

	info, err := storage.HeadBlob(d.Ctx, reponame, digest, "docker", "v2")
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); ok {
			CtxErrorWrap(d.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to head blob of '%s:%s'.", reponame, digest))
		} else {
			CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to head blob of '%s:%s'.", reponame, digest))
		}
		return
	}
	head := make(map[string]string)
	head["Content-Type"] = "application/octec-stream"
	head["Content-Length"] = fmt.Sprint(info.Size())
	CtxSuccessWrap(d.Ctx, http.StatusOK, "ok", head)
}

// GetBlob gets the blob of a certain digest
func (d *DockerV2Blob) GetBlob() {
	reponame := d.Ctx.Input.Param(":splat")
	digest := d.Ctx.Input.Param(":digest")
	logs.Debug("GetBlob of '%s:%s'.", reponame, digest)

	data, err := storage.GetBlob(d.Ctx, reponame, digest, "docker", "v2")
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); ok {
			CtxErrorWrap(d.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to get blob of '%s:%s'.", reponame, digest))
		} else {
			CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get blob of '%s:%s'.", reponame, digest))
		}
		return
	}
	CtxDataWrap(d.Ctx, http.StatusOK, data, nil)
}

// DeleteBlob deletes the blob of a certain digest
func (d *DockerV2Blob) DeleteBlob() {
	reponame := d.Ctx.Input.Param(":splat")
	digest := d.Ctx.Input.Param(":digest")
	logs.Debug("DeleteBlob of '%s:%s'.", reponame, digest)

	err := storage.DeleteBlob(d.Ctx, reponame, digest, "docker", "v2")
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); ok {
			CtxErrorWrap(d.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to delete blob of '%s:%s'.", reponame, digest))
		} else {
			CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete blob of '%s:%s'.", reponame, digest))
		}
		return
	}
	CtxSuccessWrap(d.Ctx, http.StatusOK, fmt.Sprintf("Succeed in deleting blob of '%s:%s'.", reponame, digest), nil)
}

// PostBlob posts the blob and gets a uuid
func (d *DockerV2Blob) PostBlob() {
	reponame := d.Ctx.Input.Param(":splat")
	mount := d.Ctx.Input.Query("mount")
	logs.Debug("PostBlob of '%s:[%s]'.", reponame, mount)

	uuid, err := session.New(*d.Ctx, mount)
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to create session to upload blob to '%s'.", reponame))
		return
	}
	header := make(map[string]string)
	header["Docker-Upload-UUID"] = uuid
	header["Range"] = "0-0"
	header["Content-Length"] = "0"
	header["Location"] = fmt.Sprintf("%s%s", d.Ctx.Input.URL(), uuid)
	CtxSuccessWrap(d.Ctx, http.StatusAccepted, "ok", header)
}

// PatchBlob starts to upload the blob data
func (d *DockerV2Blob) PatchBlob() {
	var uuid string
	reponame := d.Ctx.Input.Param(":splat")
	mount := d.Ctx.Input.Query("mount")
	if mount == "" {
		uuid = d.Ctx.Input.Param(":uuid")
	} else {
		uuid = mount
	}
	// FIXME: Warn: for security reason, we should not output the uuid
	logs.Debug("PatchBlob of '%s:%s'.", reponame, uuid)
	_, err := session.Get(*d.Ctx, uuid)
	if err != nil {
		// TODO: not found error
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get session in patching blob to '%s'.", reponame))
		return
	}

	data := d.Ctx.Input.CopyBody(utils.MaxSize)
	err = session.PutCache(*d.Ctx, uuid, data)
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to cache the patched blob to '%s'.", reponame))
		return
	}

	header := make(map[string]string)
	header["Docker-Upload-UUID"] = uuid
	header["Content-Length"] = "0"
	header["Location"] = fmt.Sprintf("%s", d.Ctx.Input.URL())
	header["Range"] = fmt.Sprintf("0-%v", len(data)-1)

	CtxSuccessWrap(d.Ctx, http.StatusNoContent, fmt.Sprintf("Succeed in patch blob to '%s'.", reponame), header)
}

// PutBlob marks the blob uploading status to done
func (d *DockerV2Blob) PutBlob() {
	var uuid string
	reponame := d.Ctx.Input.Param(":splat")
	mount := d.Ctx.Input.Query("mount")
	if mount == "" {
		uuid = d.Ctx.Input.Param(":uuid")
	} else {
		uuid = mount
	}
	// FIXME: Warn: for security reason, we should not output the uuid
	logs.Debug("PutBlob of '%s:%s'.", reponame, uuid)

	_, err := session.Get(*d.Ctx, uuid)
	if err != nil {
		// TODO: not found error
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get session in putting blob to '%s'.", reponame))
		return
	}

	data, err := session.GetCache(*d.Ctx, uuid)
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get cached session data in putting blob to '%s'.", reponame))
		return
	}

	err = storage.PutBlob(d.Ctx, reponame, "docker", "v2", data)
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put blob of '%s'.", reponame))
		return
	}

	err = session.Release(*d.Ctx, uuid)
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to release session in putting blob to '%s'.", reponame))
		return
	}

	digest := utils.GetDigest("sha256", data)
	header := make(map[string]string)
	header["Content-Length"] = "0"
	header["Content-Range"] = fmt.Sprintf("0-%v", len(data)-1)
	header["Docker-Content-Digest"] = digest
	CtxSuccessWrap(d.Ctx, http.StatusNoContent, fmt.Sprintf("Succeed in putting blob of '%s'.", reponame), nil)
}

// DockerV2Repo defines the repo operations
type DockerV2Repo struct {
	beego.Controller
}

// GetRepoList returns the repo list
func (d *DockerV2Repo) GetRepoList() {
	logs.Debug("GetRepoList")

	repos, err := models.QueryReposList()
	if err != nil {
		CtxErrorWrap(d.Ctx, http.StatusInternalServerError, err, "Fail to get repos list")
		return
	}

	type cataLog struct {
		Repositories []string
	}
	var c cataLog
	c.Repositories = repos
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	CtxSuccessWrap(d.Ctx, http.StatusOK, c, header)
}
