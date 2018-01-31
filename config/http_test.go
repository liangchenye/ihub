package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPValid(t *testing.T) {
	cases := []struct {
		config   HTTPConfig
		expected bool
	}{
		{HTTPConfig{Schema: "https"}, true},
		{HTTPConfig{Schema: "https"}, true},
		{HTTPConfig{Schema: "invalid"}, false},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.config.Valid() == nil)
	}
}

func TestInitHTTP(t *testing.T) {
	cases := []struct {
		config   HTTPConfig
		expected bool
	}{
		{HTTPConfig{Schema: "https", Port: 0}, true},
		{HTTPConfig{Schema: "https", Port: 443}, true},
		{HTTPConfig{Schema: "http", Port: 0}, true},
		{HTTPConfig{Schema: "http", Port: 80}, true},
		{HTTPConfig{Schema: "invalid", Port: 80}, false},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, InitHTTP(c.config) == nil)
	}
}

func TestGetDisplayAddr(t *testing.T) {
	sysConfig.Server.Schema = "http"
	sysConfig.Server.Port = 80

	cases := []struct {
		website  string
		expected string
	}{
		{"", ""},
		{"localhost", "http://localhost:80"},
	}

	for _, c := range cases {
		sysConfig.Server.Website = c.website
		assert.Equal(t, c.expected, GetDisplayAddr())
	}
}
