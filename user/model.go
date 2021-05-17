package user

import (
	"errors"
	"html/template"
	"net/mail"
	"unicode"

	"github.com/yuhsuan105/go_chatroom/common"
)

type User struct {
	Id       int `gorm:"primaryKey"`
	Name     string
	Email    string
	Password string
}

type templForm struct {
	CsrfField template.HTML
	ErrMsg    string
}

const (
	PWD_MIN_LENGTH      = 8
	NICKNAME_MIN_LENGTH = 4
)

func registerValidate(email, pwd, retype, nickname string) error {
	if err := emailValidate(email); err != nil {
		return err
	}
	if pwd != retype {
		return errors.New("password and retype password not equal")
	}
	if !passwordValidate(pwd, PWD_MIN_LENGTH) {
		return errors.New("password wrong format")
	}
	if err := nicknameValidate(nickname, NICKNAME_MIN_LENGTH); err != nil {
		return err
	}

	return nil
}

func emailValidate(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}

	var count int64

	db := common.GetDatabase()
	db.Model(&User{}).Where("email = ?", email).Count(&count)

	if count > 0 {
		return errors.New("this email have already registered")
	}

	return nil
}

func nicknameValidate(name string, minLen int) error {
	r := []rune(name)
	if len(r) < minLen {
		return errors.New("nickname length too short")
	}

	var count int64

	db := common.GetDatabase()
	db.Table("users").Where("name = ?", name).Count(&count)

	if count > 0 {
		return errors.New("this nickname have been used")
	}
	return nil
}

func passwordValidate(pwd string, minLen int) bool {
	var f_num, f_upper, f_lower bool

	if len(pwd) < minLen {
		return false
	}

	for _, c := range pwd {
		switch {
		case unicode.IsNumber(c):
			f_num = true
		case unicode.IsUpper(c):
			f_upper = true
		case unicode.IsLower(c):
			f_lower = true
		default:
			return false
		}
	}
	return f_num && f_upper && f_lower
}
