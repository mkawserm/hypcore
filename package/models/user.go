package models

import (
	"github.com/mkawserm/hypcore/package/constants"
	"github.com/mkawserm/hypcore/package/crypto/hasher"
	"time"
)

type User struct {
	Pk        string
	Password  string
	Group     string
	CreatedAt int64
	UpdatedAt int64
}

func NewSuperUser(userName string, password string) *User {
	u := &User{Pk: userName}
	u.SetGroup(constants.SuperGroup)
	u.SetPassword(password)
	u.CreatedAt = time.Now().UnixNano()
	u.UpdatedAt = time.Now().UnixNano()
	return u
}

func NewServiceUser(userName string, password string) *User {
	u := &User{Pk: userName}
	u.SetGroup(constants.ServiceGroup)
	u.SetPassword(password)

	u.CreatedAt = time.Now().UnixNano()
	u.UpdatedAt = time.Now().UnixNano()

	return u
}

func NewNormalUser(userName string, password string) *User {
	u := &User{Pk: userName}
	u.SetGroup(constants.NormalGroup)
	u.SetPassword(password)

	u.CreatedAt = time.Now().UnixNano()
	u.UpdatedAt = time.Now().UnixNano()

	return u
}

func (u *User) IsSuperUser() bool {
	return u.Group == constants.SuperGroup
}

func (u *User) IsServiceUser() bool {
	return u.Group == constants.ServiceGroup
}

func (u *User) IsNormalUser() bool {
	return u.Group == constants.NormalGroup
}

func (u *User) GetGroup() string {
	return u.Group
}

func (u *User) IsPasswordValid(rawPassword string) bool {
	b, _ := hasher.CheckPassword(rawPassword, u.Password)
	return b
}

func (u *User) SetPassword(rawPassword string) bool {
	var err error
	u.Password, err = hasher.MakePassword(rawPassword, "", "default")

	if err == nil {
		return true
	} else {
		return false
	}
}

func (u *User) SetGroup(group string) {
	u.Group = group
}
