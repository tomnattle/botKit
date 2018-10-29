package errorHandler

func config() []*errConfig {
	return []*errConfig{
		// Client error
		&errConfig{code: 20034001, msg: "参数错误"},
		// Server error
		&errConfig{code: 20035001, msg: "系统错误"},
	}
}
