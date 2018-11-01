package errorHandler

func errCodeConfig() []*errConfig {
	return []*errConfig{
		// Client error
		&errConfig{code: 4001, msg: "参数错误"},
		&errConfig{code: 4002, msg: "验证码已存在"},
		// Server error
		&errConfig{code: 5001, msg: "系统错误"},
	}
}
