package memory

import (
	"testing"

	"github.com/astaxie/beego/context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	var m Memory
	var ctx context.Context
	testData := []byte("test data")

	// Not inited
	_, err := m.New(ctx, "")
	assert.NotNil(t, err)
	_, err = m.Get(ctx, "")
	assert.NotNil(t, err)
	err = m.Release(ctx, "")
	assert.NotNil(t, err)
	_, err = m.GetCache(ctx, "")
	assert.NotNil(t, err)
	err = m.PutCache(ctx, "", testData)
	assert.NotNil(t, err)
	err = m.GC()
	assert.NotNil(t, err)

	// Inited
	err = m.Init(nil)
	assert.Nil(t, err)

	id, err := m.New(ctx, "")
	assert.Nil(t, err)

	// id is already exists
	_, err = m.New(ctx, id)
	assert.NotNil(t, err)

	// id is new
	wantedID := uuid.NewV4().String()
	id, err = m.New(ctx, wantedID)
	assert.Nil(t, err)
	assert.Equal(t, wantedID, id)

	_, err = m.Get(ctx, id)
	assert.Nil(t, err)
	_, err = m.Get(ctx, id+"-invalid")
	assert.NotNil(t, err)

	err = m.PutCache(ctx, id, testData)
	assert.Nil(t, err)
	_, err = m.GetCache(ctx, id)
	assert.Nil(t, err)
	_, err = m.GetCache(ctx, id+"-invalid")
	assert.NotNil(t, err)

	err = m.Release(ctx, id)
	assert.Nil(t, err)
	// Cannot get it after release
	_, err = m.Get(ctx, id)
	assert.NotNil(t, err)
	// Cannot get cache after release
	_, err = m.GetCache(ctx, id)
	assert.NotNil(t, err)
	// Cannot release after release
	err = m.Release(ctx, id)
	assert.NotNil(t, err)

	err = m.GC()
	assert.Nil(t, err)
}
