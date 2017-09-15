package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/controllers"
)

const (
	dockerV2Prefix = "/v2"
)

func init() {
	if err := RegisterRouter(dockerV2Prefix, DockerV2NameSpace()); err != nil {
		logs.Error("Failed to register router: '%s'.", dockerV2Prefix)
	} else {
		logs.Debug("Register router '%s' registered.", dockerV2Prefix)
	}
}

// DockerV2NameSpace defines the docker v2 router
func DockerV2NameSpace() *beego.Namespace {
	ns := beego.NewNamespace(dockerV2Prefix,
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),
		beego.NSRouter("/", &controllers.DockerV2Ping{}, "get:Ping"),
		beego.NSRouter("/_catalog", &controllers.DockerV2Repo{}, "get:GetRepoList"),
		beego.NSRouter("/*/tags/list", &controllers.DockerV2Tag{}, "get:GetTagsList"),
		//FIXME: delete the upload blob/ get the upload blob/
		beego.NSRouter("/*/blobs/uploads/?:uuid", &controllers.DockerV2Blob{}, "post:PostBlob;patch:PatchBlob;put:PutBlob"),
		beego.NSRouter("/*/blobs/:digest", &controllers.DockerV2Blob{}, "head:HeadBlob;get:GetBlob;delete:DeleteBlob"),
		beego.NSRouter("/*/manifests/:tags", &controllers.DockerV2Manifest{}, "get:GetManifest;put:PutManifest;delete:DeleteManifest"),
	)

	return ns
}
