package models

import (
	"time"

	"github.com/ifchange/botKit/xorm"
)

// Note: struct field string with digit will not split by '_', so 'Axis0' will be 'axis0' with xorm, not 'axis_0'.

type Questions struct {
	xorm.Base `xorm:"extends"`

	Id      int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Type    int    `json:"type" xorm:"not null default 0 comment('0:单选题 1:封闭主观题 2:非封闭主观题') index(idx_type_axis_is_deleted) INT(11)"`
	Axis0   int    `json:"axis0" xorm:"not null default 0 index(idx_type_axis_is_deleted) INT(11)"`
	Axis1   int    `json:"axis1" xorm:"not null default 0 index(idx_type_axis_is_deleted) INT(11)"`
	Axis2   int    `json:"axis2" xorm:"not null default 0 index(idx_type_axis_is_deleted) INT(11)"`
	Text    string `json:"text" xorm:"not null TEXT"`
	Options string `json:"options" xorm:"not null default '' comment('optionID,optionID') VARCHAR(255)"`
}

func (m *Questions) GetId() (val int) {
	if m == nil {
		return
	}
	return m.Id
}

func (m *Questions) GetType() (val int) {
	if m == nil {
		return
	}
	return m.Type
}

func (m *Questions) GetAxis0() (val int) {
	if m == nil {
		return
	}
	return m.Axis0
}

func (m *Questions) GetAxis1() (val int) {
	if m == nil {
		return
	}
	return m.Axis1
}

func (m *Questions) GetAxis2() (val int) {
	if m == nil {
		return
	}
	return m.Axis2
}

func (m *Questions) GetText() (val string) {
	if m == nil {
		return
	}
	return m.Text
}

func (m *Questions) GetOptions() (val string) {
	if m == nil {
		return
	}
	return m.Options
}

func (m *Questions) GetIsDeleted() (val int) {
	if m == nil {
		return
	}
	return m.IsDeleted
}

func (m *Questions) GetUpdatedAt() (val time.Time) {
	if m == nil {
		return
	}
	return m.UpdatedAt
}

func (m *Questions) GetCreatedAt() (val time.Time) {
	if m == nil {
		return
	}
	return m.CreatedAt
}

func (m *Questions) TableName() string {
	return "questions"
}

func CreateQuestions(obj *Questions) (int64, error) {
	return xorm.ORM().Insert(obj)
}

func UpdateQuestions(id int, obj *Questions) (int64, error) {
	return xorm.ORM().Id(id).Update(obj)
}

func DeleteQuestions(id int, obj *Questions) (int64, error) {
	return xorm.ORM().Id(id).Delete(obj)
}

func SoftDeleteQuestionsByID(id int, obj *Questions) (int64, error) {
	obj.IsDeleted = 1
	return xorm.ORM().Id(id).Update(obj)
}

func GetQuestionsByID(id int64, obj *Questions) error {
	has, err := xorm.ORM().Id(id).Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return xorm.ErrNotExist
	}
	return nil
}

func QuestionsSearch(cond *xorm.Conditions) (ts []Questions, err error) {
	if cond == nil {
		cond = xorm.NewConditions()
	}

	query, args := cond.Parse()

	if cond.Limit != 0 || cond.Offset != 0 {
		err = xorm.ORM().Where(query, args...).Limit(cond.Limit, cond.Offset).Find(&ts)
	} else {
		err = xorm.ORM().Where(query, args...).Find(&ts)
	}

	if err != nil {
		return
	}

	return
}
