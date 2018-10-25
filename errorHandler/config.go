package errorHandler

func configError() []*errConfig {
	return []*errConfig{
		// Client error
		&errConfig{code: 20044001, msg: "参数错误"},
		// Server error
		&errConfig{code: 20045001, msg: "系统错误"},
	}
}
