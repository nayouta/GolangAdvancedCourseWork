//包名称命名为小写
package handleSQLError

import (
	"database/sql"
	xerrors "github.com/pkg/errors"
)

type SQLOperationError struct {
	error
}

func (*SQLOperationError) Error() string {
	return "SQLOperation error occured"
}

type ValueFromTable struct {

}

func (*ValueFromTable) ResultDisplay(){

}

//对数据库操作API定义中的返回值一般为 (数据, error)的形式
func sqlOperation() (ValueFromTable, error) {
	return ValueFromTable{}, sql.ErrNoRows
}

func Test0() error {
	if value, err := sqlOperation(); err != nil {
		//使用pkg/errors提供的wrap方法进行包装
		return xerrors.Wrapf(err, "Test0 SQLOperation failed ") //抛给上层的同时，记录当前信息
	}else{
		//使用value
		value.ResultDisplay()
	}
	return nil
}

//返回自定义Error的指针
func Test1() *SQLOperationError {
	return &SQLOperationError{
		Test0(),
	}
}
