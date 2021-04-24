package main

import (
	"errors"
	"fmt"
	xerrors "github.com/pkg/errors"
)

var errMy = errors.New("My");

func main()  {
	err := test0()
	fmt.Printf("main:%+v \n", err)


}

func test0 () error{
	return xerrors.Wrapf(errMy,"Test0 failed")
}
