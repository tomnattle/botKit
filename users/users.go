package users

import (
	"fmt"
)

const (
	ConstManagerUser int = 1
	ConstFreeUser    int = 2
)

func GetUserType(userID int) (userType int, err error) {
	if userID <= 0 {
		err = fmt.Errorf("userID is not allowed")
	}
	if userID%2 == 0 {
		userType = ConstFreeUser
		return
	}
	userType = ConstManagerUser
	return
}
