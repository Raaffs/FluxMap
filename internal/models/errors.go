package models

import "errors"
var(
	ErrAlreadyExist=errors.New("Already exist")
	ErrInvalidCredential=errors.New("Invalid Credentials")
	ErrRecordNotFound = errors.New("record not found")
	// ErrEditConflict is returned when a there is a data race, and we have an edit conflict.
	ErrEditConflict = errors.New("edit conflict")
	
)
