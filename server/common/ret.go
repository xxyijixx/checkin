package common

type RetMessage[T any] struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func Success[T any]() *RetMessage[T] {
	return &RetMessage[T]{
		Ret: 1,
		Msg: "success",
	}
}

func SuccessWithData[T any](data T) *RetMessage[T] {
	return &RetMessage[T]{
		Ret:  1,
		Msg:  "success",
		Data: data,
	}
}

func Error[T any](msg string) *RetMessage[T] {
	return &RetMessage[T]{
		Ret: 0,
		Msg: msg,
	}
}

func ErrorWithData[T any](msg string, data T) *RetMessage[T] {
	return &RetMessage[T]{
		Ret:  0,
		Msg:  msg,
		Data: data,
	}
}
