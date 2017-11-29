package storage

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

const octPrefix = "oct"

var (
	ok   []byte
	fail []byte
)

func initData() {
	var ctx context.Context
	ok, _ = Driver().GetContent(ctx, fmt.Sprintf("%s-icons/ok", octPrefix))
	fail, _ = Driver().GetContent(ctx, fmt.Sprintf("%s-icons/fail", octPrefix))
}

func checkData(data []byte) []byte {
	if len(ok) == 0 {
		initData()
	}
	content := string(data)
	strs := strings.Split(content, "\n")
	for end := len(strs) - 1; end >= 0; end-- {
		if strings.HasPrefix(strs[end], "Result: PASS") {
			return ok
		}
	}

	return fail
}

// GetOctStatus returns the status of a repo
func GetOctStatus(ctx *context.Context, repo string) ([]byte, error) {
	storagePath := fmt.Sprintf("%s/%s/status", octPrefix, repo)
	logs.Debug("Get '%s'.", storagePath)

	return Driver().GetContent(*ctx, storagePath)
}

// GetOctImage returns the image of a repo
func GetOctImage(ctx *context.Context, repo string) ([]byte, error) {
	storagePath := fmt.Sprintf("%s/%s/image", octPrefix, repo)
	logs.Debug("Get '%s'.", storagePath)

	return Driver().GetContent(*ctx, storagePath)
}

// AddOctOutput writes the output to storage
func AddOctOutput(ctx *context.Context, repo string, data []byte) error {
	storageDir := fmt.Sprintf("%s/%s", octPrefix, repo)
	logs.Debug("AddOctOutput '%s'.", storageDir)

	imageData := checkData(data)
	err := Driver().PutContent(*ctx, filepath.Join(storageDir, "image"), imageData)
	if err != nil {
		return err
	}
	//FIXME: rollback to remove 'image'
	return Driver().PutContent(*ctx, filepath.Join(storageDir, "status"), data)
}
