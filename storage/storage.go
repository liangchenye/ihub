package storage

import (
	"errors"

	beegoContext "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	dockerContext "github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/storage/driver"
	"github.com/docker/distribution/registry/storage/driver/factory"

	"github.com/isula/ihub/config"
	"github.com/isula/ihub/health"
)

// Health provides interface for storage health operations
type Health struct {
}

// GetStatus returns the storage health status
func (sh *Health) GetStatus() (string, string) {
	return "", ""
}

// HealthRegist regists the storage health driver
func HealthRegist() error {
	health.RegisterHealth("storage", &Health{})
	return nil
}

var (
	sysDriver driver.StorageDriver
)

func init() {
	sysDriver = nil
}

// TODO: more logs
func loadDriver(cfg config.StorageConfig) (driver.StorageDriver, error) {
	for n, paras := range cfg {
		logs.Debug("Find storage driver for: %s, %v", n, paras)
		d, err := factory.Create(n, paras)
		if err == nil {
			return d, nil
		}
	}

	return nil, errors.New("Fail to get a suitable storage driver")
}

// InitStorage setups the storage bankends from the config file
func InitStorage(cfg config.StorageConfig) error {
	var err error
	sysDriver, err = loadDriver(cfg)
	// TODO: we should check the healthy status at the beginning

	return err
}

// Driver returns the storage driver
func Driver() driver.StorageDriver {
	cfg := config.GetConfig()
	if cfg.StorageLoad == "static" && sysDriver != nil {
		return sysDriver
	}

	var err error
	sysDriver, err = loadDriver(cfg.Storage)
	if err != nil {
		logs.Debug(err)
		panic("Failed to load driver")
	}

	return sysDriver
}

// BC2DC converts beego context to docker context
func BC2DC(bctx *beegoContext.Context) *dockerContext.Context {
	//TODO
	var ctx dockerContext.Context
	return &ctx
}
