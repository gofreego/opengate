package customerrors

import (
	"github.com/gofreego/goutils/customerrors"
)

var (
	ErrNoJWTToken = customerrors.BAD_REQUEST_ERROR("no jwt token")
)
