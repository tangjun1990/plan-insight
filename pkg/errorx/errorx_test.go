package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(errors.New("无效的参数uid"), CodeInternalServer, map[string]string{
		"aaa": "bbb",
		"ccc": "ddd",
	})
	fmt.Println(err.Err)
	fmt.Println(err.Error())
	fmt.Println(err.Code)
	fmt.Println(err.Meta)
	fmt.Println(err.JSON())
}

func TestNewError(t *testing.T) {
	err := NewArgument("无效的参数uid")
	fmt.Println(err.Err)
	fmt.Println(err.Error())
	fmt.Println(err.Code)
	fmt.Println(err.Meta)
	fmt.Println(err.JSON())
}

func TestNewErrore(t *testing.T) {
	e := errors.New("无效的参数uid")

	err := NewArgumente(e)

	fmt.Println(err.Err)
	fmt.Println(err.Error())
	fmt.Println(err.Code)
	fmt.Println(err.Meta)
	fmt.Println(err.JSON())
}

func TestNewErrorf(t *testing.T) {
	e := errors.New("无效的参数uid")

	err := NewArgumentf("err： %v, uid: %d", e, 123456)

	fmt.Println(err.Err)
	fmt.Println(err.Error())
	fmt.Println(err.Code)
	fmt.Println(err.Meta)
	fmt.Println(err.JSON())
}
