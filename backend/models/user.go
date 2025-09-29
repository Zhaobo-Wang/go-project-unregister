package models

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"` // "-" 表示不在 JSON 中显示
	Todos    []Todo `json:"todos" gorm:"foreignKey:UserID"`
}

// 更新 Todo 模型以包含用户关联
func UpdateTodoModel() {
	// 在 models/todo.go 中添加 UserID 字段
}

// HashPassword 加密用户密码
/**
参数 (Parameters)
接收者参数 (Receiver):
u *User: 指向 User 结构体的指针，表示这是 User 结构体的方法
使用指针接收者是因为需要修改 User 实例的 Password 字段
输入参数:
password string: 用户的原始明文密码
类型: 字符串 (string)
返回值 (Return)
error:
如果密码加密过程中出现错误，返回相应的错误信息
如果加密成功，返回 nil
类型: error 接口
函数内部变量
**/
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证用户密码
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}