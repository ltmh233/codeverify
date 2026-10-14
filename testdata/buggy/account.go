package buggy

import "errors"

var ErrInsufficientBalance = errors.New("not enough balance to transfer!")

type Account struct {
	Balance int64
}

type User struct {
	Balance int64
}

func (a *Account) Save() {

}

func (u *User) Save() {

}
