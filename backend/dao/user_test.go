package dao

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserInfo_Insert(t *testing.T) {
	info := &UserInfo{
		Name:   "kkk",
		Addr:   "1234567890",
		Gender: 1,
		Memo:   "test user",
	}
	err := User.Insert(info)
	assert.Nil(t, err)

	newInfo := User.GetByAddr(info.Addr)
	assert.Equal(t, info.Name, newInfo.Name)

	_, err = User.DB().Exec(fmt.Sprintf("TRUNCATE TABLE %s", User.TableName()))
	assert.Nil(t, err)
}

func TestUserInfo_Update(t *testing.T) {
	info := &UserInfo{
		Name:   "kkk",
		Addr:   "1234567890",
		Gender: 1,
		Memo:   "test user",
	}
	ast := assert.New(t)

	err := User.Insert(info)
	ast.Nil(err)

	newInfo := User.GetByAddr(info.Addr)
	newInfo.Name = "xxx"
	err = User.Update(newInfo.ID, newInfo)
	ast.Nil(err)

	newInfo = User.GetByAddr(newInfo.Addr)
	ast.NotEqual(info.Name, newInfo.Name)

	err = User.Delete(newInfo.ID)
	ast.Nil(err)

	newInfo = User.GetByAddr(newInfo.Addr)
	ast.Equal(newInfo.Name, "")
}
