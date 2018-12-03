package {{.Models}}

import (
{{$ilen := len .Imports}}
{{if gt $ilen 0}}
{{range .Imports}}"{{.}}"{{end}}
{{end}}
    "github.com/ifchange/botKit/xorm"
)


{{range .Tables}}
type {{Mapper .Name}} struct {
	xorm.Base `xorm:"extemds"`
{{$table := .}}
{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{if eq $col.Name "id" "is_deleted" "updated_at" "created_at"}}{{else}}{{Mapper $col.Name}}	{{Type $col}} {{Tag $table $col}}{{end}}
{{end}}
}

{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}
func (m *{{Mapper $table.Name}}) Get{{Mapper $col.Name}}() (val {{Type $col}}) {
    if m == nil {
        return
    }
    return m.{{Mapper $col.Name}}
}
{{end}}

func (m *{{Mapper .Name}}) TableName() string {
	return "{{.Name}}"
}

func Create{{Mapper .Name}}(obj *{{Mapper .Name}}) (int64, error) {
	return xorm.ORM().Insert(obj)
}

func Update{{Mapper .Name}}(obj *{{Mapper .Name}}) (int64, error) {
	return xorm.ORM().Update(obj)
}

func Delete{{Mapper .Name}}(id int, obj *{{Mapper .Name}}) (int64, error) {
	return xorm.ORM().Id(id).Delete(obj)
}

func SoftDelete{{Mapper .Name}}ByID(id int, obj *{{Mapper .Name}})(int64, error) {
	obj.IsDeleted = 1
	return xorm.ORM().Id(id).Update(obj)
}

func Get{{Mapper .Name}}ByID(id int64, obj *{{Mapper .Name}}) error {
	has, err := xorm.ORM().Id(id).Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return xorm.ErrNotExist
	}
	return nil
}

func {{Mapper .Name}}Search(cond *xorm.Conditions) (ts []{{Mapper .Name}}, err error) {
    if cond == nil {
        cond = xorm.NewConditions()
    }

	query, args := cond.Parse()

	err = xorm.ORM().Where(query, args).Find(&ts)
	if err != nil {
		return
	}

	return
}

{{end}}
