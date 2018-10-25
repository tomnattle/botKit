package util

func Max(num int, nums ...int) int {
	for _, item := range nums {
		if item > num {
			num = item
		}
	}
	return num
}

func Min(num int, nums ...int) int {
	for _, item := range nums {
		if item < num {
			num = item
		}
	}
	return num
}
