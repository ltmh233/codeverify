package buggy

import "errors"

var ErrInsufficientBalance = errors.New("not enough balance to transfer!")

type Account struct {
	Balance int64
}
