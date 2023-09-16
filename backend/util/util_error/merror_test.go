package util_error

import (
	"fmt"
	"github.com/pkg/errors"
)

func wrapErr() error {
	err := errors.New("error in wrap")
	err = Wrap(err, "wrap")
	return err
}
func result1() {
	err := wrapErr()
	fmt.Printf("result err:%s", err)
}
func ExampleWrap() {
	result1()
	// OutPut: result err:==> TLDrivingInformation/util/merror.wrapErr()@10: wrap: error in wrap
}

func withErr1() error {
	err := wrapErr()
	err = WithMessage(err, "with err1")
	return err
}
func result2() {
	err := withErr1()
	fmt.Printf("reuslt err:%s", err)
}
func ExampleWithMessage() {
	result2()
	// OutPut:
	// reuslt err:==> TLDrivingInformation/util/merror.withErr1()@24: with err1: ==> TLDrivingInformation/util/merror.wrapErr()@10: wrap: error in wrap
}
