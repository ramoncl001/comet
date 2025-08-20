package rest

type Response struct {
	Status int
	Data   interface{}
}

func Ok[T any](data T) Response {
	return Response{
		Status: 200,
		Data:   data,
	}
}

func Error[T any](data T) Response {
	return Response{
		Status: 500,
		Data:   data,
	}
}

func NotFound() Response {
	return Response{
		Status: 404,
		Data:   "resource not found",
	}
}

func BadRequest[T any](data T) Response {
	return Response{
		Status: 400,
		Data:   data,
	}
}

func Unauthorized() Response {
	return Response{
		Status: 401,
		Data:   "Unauthorized",
	}
}
