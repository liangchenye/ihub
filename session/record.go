package session

import (
	"github.com/astaxie/beego/context"
)

// Record provides interface to manage a session detailed info,
// such like expiring data
type Record struct {
}

// Match checkes if the request is received from the expected user
// TODO
func (r *Record) Match(ctx context.Context) error {
	return nil
}

// Expired checkes if a session is out of data
func (r *Record) Expired() bool {
	return false
}

// NewRecordFromContext creates a record automaticlly from a http context
func NewRecordFromContext(ctx context.Context) Record {
	var r Record
	return r
}
