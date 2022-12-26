package messages

import "github.com/pkg/errors"

var ErrInvalidData = errors.New("invalid data")
var ErrLoginUsed = errors.New("login already registered")
var ErrAccountUsed = errors.New("account already registered")
var ErrInvalidToken = errors.New("token is not available")
var ErrIncorrectPassword = errors.New("incorrect password")
var ErrExistsAccount = errors.New("such account does not exist")
var ErrExpiredToken = errors.New("token is expired")
var ErrAccountBelongUser = errors.New("incorrect account entered or does not belong to the user")
var ErrExistsType = errors.New("this type of operation is not registered")
