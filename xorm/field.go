package xorm

import (
	"time"
)

type Base struct {
	IsDeleted int       `json:"is_deleted" xorm:"not null default 0 index(idx_type_axis_is_deleted) TINYINT(4)"`
	UpdatedAt time.Time `json:"updated_at" xorm:"default 'CURRENT_TIMESTAMP' updated TIMESTAMP"`
	CreatedAt time.Time `json:"created_at" xorm:"default 'CURRENT_TIMESTAMP' created TIMESTAMP"`
}
