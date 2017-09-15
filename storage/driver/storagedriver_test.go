package driver

import (
	"errors"
	"io"
	"testing"

	"github.com/astaxie/beego/context"
	"github.com/stretchr/testify/assert"
)

type MockDriver struct {
	name string
}

func (d *MockDriver) Init(parameters map[string]interface{}) error {
	return nil
}

func (d *MockDriver) Valid(parameters map[string]interface{}) error {
	if parameters != nil {
		v := parameters["name"].(string)
		if v != "" {
			return nil
		}
	}

	return errors.New("missing 'name' data")
}

// Implement the storagedriver.StorageDriver interface

func (d *MockDriver) Name() string {
	return "mock"
}

// GetContent retrieves the content stored at "path" as a []byte.
func (d *MockDriver) GetContent(ctx context.Context, path string) ([]byte, error) {
	return nil, nil
}

// PutContent stores the []byte content at a location designated by "path".
func (d *MockDriver) PutContent(ctx context.Context, subPath string, contents []byte) error {
	return nil
}

// Reader retrieves an io.ReadCloser for the content stored at "path" with a
// given byte offset.
func (d *MockDriver) Reader(ctx context.Context, path string, offset int64) (io.ReadCloser, error) {
	return nil, nil
}

func (d *MockDriver) Writer(ctx context.Context, subPath string, append bool) (FileWriter, error) {
	return nil, nil
}

// Stat retrieves the FileInfo for the given path, including the current size
// in bytes and the creation time.
func (d *MockDriver) Stat(ctx context.Context, subPath string) (FileInfo, error) {
	return nil, nil
}

// List returns a list of the objects that are direct descendants of the given
// path.
func (d *MockDriver) List(ctx context.Context, subPath string) ([]string, error) {
	return nil, nil
}

// Move moves an object stored at sourcePath to destPath, removing the original
// object.
func (d *MockDriver) Move(ctx context.Context, sourcePath string, destPath string) error {
	return nil
}

// Delete recursively deletes all objects stored at "path" and its subpaths.
func (d *MockDriver) Delete(ctx context.Context, subPath string) error {
	return nil
}

// URLFor returns a URL which may be used to retrieve the content stored at the given path.
// May return an UnsupportedMethodErr in certain StorageDriver implementations.
func (d *MockDriver) URLFor(ctx context.Context, path string, options map[string]interface{}) (string, error) {
	return "", nil
}

func TestErrors(t *testing.T) {
	var eum ErrUnsupportedMethod
	t.Log(eum.Error())

	var pnfe PathNotFoundError
	t.Log(pnfe.Error())

	var ipe InvalidPathError
	t.Log(ipe.Error())

	var ioe InvalidOffsetError
	t.Log(ioe.Error())

	var e Error
	t.Log(e.Error())

}

func TestRegister(t *testing.T) {
	cases := []struct {
		name     string
		driver   StorageDriver
		expected bool
	}{
		{"TestRegister-A", &MockDriver{}, true},
		{"TestRegister-A", &MockDriver{}, false},
		{"", &MockDriver{}, false},
		{"TestRegister-B", nil, false},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, Register(c.name, c.driver) == nil)
	}

}

func TestFindDriver(t *testing.T) {
	Register("TestFindDriver-A", &MockDriver{})
	validParams := make(map[string]interface{})
	validParams["name"] = "/tmp"

	cases := []struct {
		name     string
		params   map[string]interface{}
		expected bool
	}{
		{"TestFindDriver-A", validParams, true},
		{"TestFindDriver-A", nil, false},
		{"TestFindDriver-B", nil, false},
	}

	for _, c := range cases {
		_, err := FindDriver(c.name, c.params)
		assert.Equal(t, c.expected, err == nil)
	}
}
