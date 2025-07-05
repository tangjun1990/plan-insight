package kgorm

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type (
	Cond struct {
		Op  string
		Val interface{}
	}
	Conds map[string]interface{}
	Ups   = map[string]interface{}
)

func assertCond(cond interface{}) Cond {
	switch v := cond.(type) {
	case Cond:
		return v
	case string:
		return Cond{"=", v}
	case bool:
		return Cond{"=", v}
	case float64:
		return Cond{"=", v}
	case float32:
		return Cond{"=", v}
	case int:
		return Cond{"=", v}
	case int64:
		return Cond{"=", v}
	case int32:
		return Cond{"=", v}
	case int16:
		return Cond{"=", v}
	case int8:
		return Cond{"=", v}
	case uint:
		return Cond{"=", v}
	case uint64:
		return Cond{"=", v}
	case uint32:
		return Cond{"=", v}
	case uint16:
		return Cond{"=", v}
	case uint8:
		return Cond{"=", v}
	case time.Duration:
		return Cond{"=", v}
	}
	condValueStr, err := cast.ToStringSliceE(cond)
	if err == nil {
		return Cond{"in", condValueStr}
	}
	condValueInt, err := cast.ToIntSliceE(cond)
	if err == nil {
		return Cond{"in", condValueInt}
	}
	log.Printf("[assertCond] unrecognized type fail,%+v\n", cond)
	return Cond{}
}

func BuildQuery(conds Conds) (sql string, binds []interface{}) {
	sql = "1=1"
	binds = make([]interface{}, 0, len(conds))
	for field, cond := range conds {
		condVal := assertCond(cond)

		// 说明有表的数据
		if strings.Contains(field, ".") {
			arr := strings.Split(field, ".")
			if len(arr) != 2 {
				return
			}
			field = "`" + arr[0] + "`.`" + arr[1] + "`"
		} else {
			field = "`" + field + "`"
		}

		switch strings.ToLower(condVal.Op) {
		case "like":
			if condVal.Val != "" {
				sql += " AND " + field + " like ?"
				condVal.Val = "%" + condVal.Val.(string) + "%"
			}
		case "%like":
			if condVal.Val != "" {
				sql += " AND " + field + " like ?"
				condVal.Val = "%" + condVal.Val.(string)
			}
		case "like%":
			if condVal.Val != "" {
				sql += " AND " + field + " like ?"
				condVal.Val = condVal.Val.(string) + "%"
			}
		case "in", "not in":
			sql += " AND " + field + condVal.Op + " (?) "
		case "between":
			sql += " AND " + field + condVal.Op + " ? AND ?"
			val := cast.ToStringSlice(condVal.Val)
			binds = append(binds, val[0], val[1])
			continue
		case "exp":
			sql += " AND " + field + " ? "
			condVal.Val = gorm.Expr(condVal.Val.(string))
		default:
			sql += " AND " + field + condVal.Op + " ? "
		}
		binds = append(binds, condVal.Val)
	}
	return
}
