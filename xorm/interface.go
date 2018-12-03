package xorm

import "time"

// Table is the database table object.
type Table interface {
	SetDeleted()
	IsDeleted() bool
	Recover()
}

// Deletable has deleted contribute and implements Tabble interface.
type Deletable struct {
	Deleted bool `json:"deleted" xorm:"deleted"`
}

// IsDeleted returns the obj is whether deleted or not.
func (d *Deletable) IsDeleted() bool {
	return d.Deleted
}

// SetDeleted sets the obj as deleted.
func (d *Deletable) SetDeleted() {
	d.Deleted = true
}

// Recover sets the obj to not deleted.
func (d *Deletable) Recover() {
	d.Deleted = false
}

// Updatable has created and updated contribute.
type Updatable struct {
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}
