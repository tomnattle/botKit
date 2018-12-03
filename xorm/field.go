package xorm

import (
	"time"
)

type Base struct {
	Id        int       `json:"id" xorm:"not null pk autoincr INT(10)"`
	IsDeleted int       `json:"is_deleted" xorm:"not null default 0 index(idx_question_id_is_deleted) TINYINT(1)"`
	UpdatedAt time.Time `json:"updated_at" xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	CreatedAt time.Time `json:"created_at" xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
}
