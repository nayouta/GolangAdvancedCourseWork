package HandleSQLError

import (
	"database/sql"
	xerrors "github.com/pkg/errors"
)


type SQLOperationError struct{
	error
}

func ( *SQLOperationError) Error() string {
	return "SQLOperation error occured"
}

func sqlOperation() (bool, error){
	return false, sql.ErrNoRows
}

func Test0 () error{
	if ok, err := sqlOperation(); ok == false{
		//使用pkg/errors提供的wrap方法进行包装
		return xerrors.Wrapf(err,"Test0 SQLOperation failed ")	//抛给上层的同时，记录当前信息
	}
	return nil
}

//返回自定义Error的指针
func Test1() *SQLOperationError{
	return &SQLOperationError{
		Test0(),
	}
}



