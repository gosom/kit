package lib_test

import "github.com/gosom/kit/lib"

var _ lib.User = &User{}

type User struct {
	IDFn    func() string
	ExtraFn func() map[string]string
}

func (u *User) GetID() string {
	if u.IDFn != nil {
		return u.IDFn()
	}
	return "123"
}

func (u *User) GetExtra() map[string]string {
	if u.ExtraFn != nil {
		return u.ExtraFn()
	}
	return map[string]string{"foo": "bar"}
}
