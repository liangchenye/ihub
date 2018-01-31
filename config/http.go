package config

import (
	"fmt"

	"github.com/astaxie/beego"
)

// HTTPConfig stores the http config item group
type HTTPConfig struct {
	// Schema: 'https' or 'http'
	Schema string `yaml:"schema,omitempty"`
	// Addr: listen address, default to empty and listen to all
	Addr string `yaml:"addr,omitempty"`
	// Port: default to 8080
	Port int `yaml:"port"`
	// KeyFile: private key
	KeyFile string `yaml:"keyfile,omitempty"`
	// CertFile: server certification file
	CertFile string `yaml:"certfile,omitempty"`
	// Website: show the ip/domain
	Website string `yaml:"website,omitempty"`
	// DisplayAdd: @schema://@website:@port
	DisplayAddr string `yaml:"displayaddr,omitempty"`
}

// Valid validates the http config
func (cfg *HTTPConfig) Valid() error {
	// TODO: more validation, for example, the file existence
	if cfg.Schema != "https" && cfg.Schema != "http" {
		return fmt.Errorf("invalid schema '%s', only 'http' and 'https' are supported ", cfg.Schema)
	}

	return nil
}

// InitHTTP turns config to Beego style
func InitHTTP(cfg HTTPConfig) error {
	if cfg.Schema == "https" {
		beego.BConfig.Listen.EnableHTTPS = true
		beego.BConfig.Listen.EnableHTTP = false
		if cfg.Port != 0 {
			beego.BConfig.Listen.HTTPSPort = cfg.Port
		}

		beego.BConfig.Listen.HTTPSKeyFile = cfg.KeyFile
		beego.BConfig.Listen.HTTPSCertFile = cfg.CertFile
	} else if cfg.Schema == "http" {
		if cfg.Port != 0 {
			beego.BConfig.Listen.HTTPPort = cfg.Port
		}
	} else {
		return fmt.Errorf("invalid schema '%s', only 'http' and 'https' are supported ", cfg.Schema)
	}

	return nil
}

// GetDisplayAddr gets the full display address of the website
func GetDisplayAddr() string {
	cfg := sysConfig.Server
	if cfg.Website != "" {
		return fmt.Sprintf("%s://%s:%d", cfg.Schema, cfg.Website, cfg.Port)
	}
	return ""
}
