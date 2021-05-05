package main

import (
	"database/sql"
	"errors"
	"fmt"
	errors2 "github.com/pkg/errors"
	"main.go/handleSQLError"
)

var errMy = errors.New("My")

func main() {
	if err := handleSQLError.Test0(); err != nil {
		//1、直接输出err信息
		fmt.Printf("main:%+v \n", err)
		//2、调用Is判断
		fmt.Printf("error.Is Result : %+v \n", errors2.Is(err, sql.ErrNoRows))
		//3、调用AS判断
		fmt.Printf("error.As Result : %+v \n", errors2.As(handleSQLError.Test1(), &sql.ErrNoRows))
		//4、查询根因
		fmt.Printf("error.Cause Is : %+v \n", errors2.Cause(err))
	}
}
