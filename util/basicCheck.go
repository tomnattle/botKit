package util

import (
	"fmt"
	"regexp"
)

func CheckName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("empty name")
	}
	return nil
}

func CheckPhone(phone string) error {
	match, err := regexp.MatchString(`^1[0-9]{10,10}$`, phone)
	if err != nil {
		return fmt.Errorf("RegexpPhone error phone:%v err:%v",
			phone, err)
	}
	if !match {
		return fmt.Errorf("RegexpPhone phone not allow %v", phone)
	}
	return nil
}
