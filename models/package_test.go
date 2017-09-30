package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryPkgListByRepoName(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		name     string
		len      int
		expected bool
	}{
		{"isula", 2, true},
		{"non-exist", 0, true},
	}

	for _, c := range cases {
		pkgs, err := QueryPkgListByRepoName(c.name)
		assert.Equal(t, c.expected, err == nil)
		assert.Equal(t, c.len, len(pkgs), "fail to match list length")
	}
}

func TestAddPkgRepo(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		name     string
		expected bool
	}{
		{"isula", true},
		{"tmpIsula", true},
	}

	for _, c := range cases {
		repo, err := AddPkgRepo(c.name)
		assert.Nil(t, err)
		assert.Equal(t, repo.Name, c.name)
	}
}

func TestQueryPkgByID(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		repoID      int64
		pkgname     string
		expectedPkg bool
		expectedErr bool
	}{
		{1, "ostree", true, true},
		{1000, "ostree", false, true},
		{1, "non-exist-pkg", false, true},
	}

	for _, c := range cases {
		pkg, err := QueryPkgByID(c.repoID, c.pkgname)
		assert.Equal(t, c.expectedPkg, pkg != nil)
		assert.Equal(t, c.expectedErr, err == nil)
	}
}

func TestQueryPkgByName(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		reponame    string
		pkgname     string
		expectedPkg bool
		expectedErr bool
	}{
		{"isula", "ostree", true, true},
		{"non-exist-repo", "ostree", false, true},
		{"isula", "non-exist-pkg", false, true},
	}

	for _, c := range cases {
		pkg, err := QueryPkgByName(c.reponame, c.pkgname)
		assert.Equal(t, c.expectedPkg, pkg != nil)
		assert.Equal(t, c.expectedErr, err == nil)
	}
}

func TestAddPkg(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		reponame string
		name     string
		size     int64
		pkgType  string
		expected bool
	}{
		{"isula", "new", 0, "", true},
		{"isula", "ostree", 0, "", true},
		{"new-repo", "ostree", 0, "", true},
	}

	for _, c := range cases {
		_, err := AddPkg(c.reponame, c.name, c.size, c.pkgType)
		assert.Equal(t, c.expected, err == nil)
	}
}

func TestPkgDownloadInc(t *testing.T) {
	if !testReady {
		return
	}

	old, _ := QueryPkgByName("euleros", "ostree")
	assert.NotNil(t, old)

	PkgDownloadInc("euleros", "ostree")

	new, _ := QueryPkgByName("euleros", "ostree")
	assert.NotNil(t, new)

	assert.Equal(t, old.Downloads+1, new.Downloads)
}
