package filesystem

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/astaxie/beego/context"
	"github.com/stretchr/testify/assert"
)

func TestInitAndValid(t *testing.T) {
	cases := []struct {
		paras    map[string]interface{}
		expected bool
	}{
		{nil, false},
		{map[string]interface{}{"invalid": "/tmp"}, false},
		{map[string]interface{}{"rootDirectory": 1234}, false},
		{map[string]interface{}{"rootDirectory": ""}, false},
		{map[string]interface{}{"rootDirectory": "/tmp"}, true},
	}
	for _, c := range cases {
		var d driver
		assert.Equal(t, c.expected, d.Valid(c.paras) == nil)
		assert.Equal(t, c.expected, d.Init(c.paras) == nil)
	}
}

func TestName(t *testing.T) {
	var d driver
	assert.Equal(t, driverName, d.Name())
}

func TestPutAndGetContent(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "test")
	defer os.RemoveAll(tmpDir)

	var ctx context.Context
	var d driver
	paras := map[string]interface{}{"rootDirectory": tmpDir}
	d.Init(paras)

	testDir := "testdir"
	testPath := testDir + "/testPath"
	testData := []byte("testdata")
	assert.Nil(t, d.PutContent(ctx, testPath, testData))
	assert.NotNil(t, d.PutContent(ctx, testDir, testData))

	_, err := d.GetContent(ctx, testDir)
	assert.NotNil(t, err)
	_, err = d.GetContent(ctx, testDir+"/nonexist")
	assert.NotNil(t, err)

	data, _ := d.GetContent(ctx, testPath)
	assert.Equal(t, testData, data)
}
