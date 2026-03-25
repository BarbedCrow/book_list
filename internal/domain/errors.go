package domain

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email already registered")
	ErrWrongPassword = errors.New("wrong password")

	ErrBookNotFound = errors.New("book not found")

	ErrAuthorNotFound = errors.New("author not found")

	ErrListNotFound = errors.New("list not found")
	ErrNotOwner     = errors.New("user is not the list owner")
)
