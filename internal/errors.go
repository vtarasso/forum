package internal

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateName      = errors.New("duplicate name")
	ErrInvalidParent      = errors.New("invalid parent")
	ErrInvalidObjectId    = errors.New("invalid reaction object")
	ErrBlank              = "This field cannot be blank"
	ErrUsedEmail          = "Email address is already in use"
	ErrEmail              = "This field must be a valid email address"
	ErrUsedName           = "Name is already in use"
	ErrValidEmail         = "This field must be a valid email address"
	ErrChoiceCategory     = "Choose one category"
	ErrMaxChars           = "This field cannot be more than 100 characters long"
	ErrPass               = "Password must contain capital letters and lowercase, numbers, special characters and must be at least 8 and not more 20 characters."
	ErrCorrectName        = "Write correct name. Username should start with an alphabet [A-Za-z] and all other characters can be alphabets, numbers or an underscore so, [A-Za-z0-9_]. The username consists of 5 to 15 characters inclusive."
)
