package config

import (
	"path/filepath"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConnection(t *testing.T) {
	cases := []struct {
		cfg  DBConfig
		conn string
		err  error
	}{
		{DBConfig{"mysql", "user", "passwd", "server", "name"}, "user:passwd@tcp(server)/name?charset=utf8", nil},
		{DBConfig{"mysql", "", "passwd", "server", "name"}, "", ErrEmptyDBUserOrPassword},
		{DBConfig{"mysql", "user", "", "server", "name"}, "", ErrEmptyDBUserOrPassword},
		{DBConfig{"mysql", "user", "passwd", "", "name"}, "", ErrEmptyDBServer},
		{DBConfig{"mysql", "user", "passwd", "server", ""}, "", ErrEmptyDBName},
	}

	for _, c := range cases {
		conn, err := c.cfg.GetConnection()
		assert.Equal(t, c.conn, conn, "Failed to get connection url: "+c.conn)
		assert.Equal(t, c.err, err)
	}
}

func TestInitConfigFromFile(t *testing.T) {
	cases := []struct {
		file     string
		expected bool
	}{
		{"default.yml", true},
		{"non-exist.yml", false},
		{"invalid.yml", false},
		{"invalidlog.yml", false},
	}

	for _, c := range cases {
		_, err := InitConfigFromFile(filepath.Join("testdata", c.file))
		assert.Equal(t, c.expected, err == nil, "Failed to load config file: "+c.file)
	}
}

func TestConfigValid(t *testing.T) {
	cases := []struct {
		file     string
		expected bool
	}{
		{"default.yml", true},
		{"invalidsql.yml", false},
		{"invalidstorage.yml", false},
		{"nosql.yml", false},
		{"invalidhttp.yml", false},
		{"invalidsession.yml", false},
	}

	for _, c := range cases {
		InitConfigFromFile(filepath.Join("testdata", c.file))
		cfg := GetConfig()
		assert.Equal(t, c.expected, cfg.Valid() == nil, "Failed to valid config file: "+c.file)
	}
}

func TestConfig(t *testing.T) {
	assert.Equal(t, sysConfig, GetConfig())
}
