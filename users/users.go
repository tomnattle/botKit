package users

import (
	"fmt"
)

const (
	ConstManagerUser int = 1
	ConstFreeUser    int = 2

	constFreeUserAutoIncrement = 1000000000
)

func GetUserType(userID int) (userType int, err error) {
	if userID <= 0 {
		err = fmt.Errorf("userID is not allowed")
	}
	if userID >= constFreeUserAutoIncrement {
		userType = ConstFreeUser
		return
	}
	userType = ConstManagerUser
	return
}
