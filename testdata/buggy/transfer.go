package buggy

// Transfer moves amount from src to dst.
// amount must be positive.
func Transfer(src, dst *Account, amount int64) error {
	if src.Balance < amount {
		return ErrInsufficientBalance
	}

	src.Balance -= amount
	dst.Balance += amount
	return nil
}
