package common

type ResponseSuccess struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponsePagination struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(data interface{}) ResponseSuccess {
	return ResponseSuccess{
		Code:    200,
		Message: "Successful Operation",
		Data:    data,
	}
}

func ErrorResponse(code int, message string) ResponseError {
	return ResponseError{
		Code:    code,
		Message: message,
	}
}

func PaginationResponse(page, perpage int, data interface{}) ResponsePagination {
	return ResponsePagination{
		Code:    200,
		Message: "Succesful Operation",
		Page:    page,
		PerPage: perpage,
		Data:    data,
	}
}
