package storage

import (
	"fmt"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/storage/driver"
	"github.com/isula/ihub/utils"
)

// ComposeBlobPath composes the blob path from the 'digest, proto, proto version'
// repo is not used
func ComposeBlobPath(repo string, digest string, proto string, protoVersion string) string {
	head, real := utils.Snap(digest)
	return fmt.Sprintf("%s/%s/blobs/%s/%s", proto, protoVersion, head, real)
}

// HeadBlob return the blob stat
// TODO we need to get user in ctx, or setting in config
func HeadBlob(ctx *context.Context, repo string, digest string, proto string, protoVersion string) (driver.FileInfo, error) {
	storagePath := ComposeBlobPath(repo, digest, proto, protoVersion)
	logs.Debug("Head '%s'.", storagePath)

	return Driver().Stat(*ctx, storagePath)
}

// GetBlob gets the blob data
// TODO we need to get user in ctx, or setting in config
func GetBlob(ctx *context.Context, repo string, digest string, proto string, protoVersion string) ([]byte, error) {
	storagePath := ComposeBlobPath(repo, digest, proto, protoVersion)
	logs.Debug("Get '%s'.", storagePath)

	return Driver().GetContent(*ctx, storagePath)
}

// PutBlob puts the blob
// TODO we need to get user in ctx, or setting in config
func PutBlob(ctx *context.Context, repo string, proto string, protoVersion string, data []byte) error {
	digest := utils.GetDigest("sha256", data)
	storagePath := ComposeBlobPath(repo, digest, proto, protoVersion)
	logs.Debug("Put '%s'.", storagePath)

	return Driver().PutContent(*ctx, storagePath, data)
}

// DeleteBlob deletes the blob
// TODO we need to get user in ctx, or setting in config
func DeleteBlob(ctx *context.Context, repo string, digest string, proto string, protoVersion string) error {
	storagePath := ComposeBlobPath(repo, digest, proto, protoVersion)
	logs.Debug("Delete '%s'.", storagePath)

	return Driver().Delete(*ctx, storagePath)
}
