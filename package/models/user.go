package models

import "github.com/mkawserm/hypcore/package/crypto/hasher"
import "github.com/mkawserm/hypcore/package/core"

type User struct {
	Pk       string
	Password string
	Group    string
}

func NewSuperUser(userName string, password string) *User {
	u := &User{Pk: userName}
	u.SetGroup(core.SuperGroup)
	u.SetPassword(password)
	return u
}

func NewServiceUser(userName string, password string) *User {
	u := &User{Pk: userName}
	u.SetGroup(core.ServiceGroup)
	u.SetPassword(password)
	return u
}

func NewNormalUser(userName string, password string) *User {
	u := &User{Pk: userName}
	u.SetGroup(core.NormalGroup)
	u.SetPassword(password)
	return u
}

func (u *User) IsSuperUser() bool {
	return u.Group == core.SuperGroup
}

func (u *User) IsServiceUser() bool {
	return u.Group == core.ServiceGroup
}

func (u *User) IsNormalUser() bool {
	return u.Group == core.NormalGroup
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
