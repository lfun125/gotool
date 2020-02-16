package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	e1 := errors.New("第一个错误")
	e2 := New(e1, 300)
	fmt.Println(e1)
	fmt.Println(e2)
	var target Error
	fmt.Println(errors.As(e2, &target))
	fmt.Println(target.Is(e1))
	errors.Is(e1, e1)
}
