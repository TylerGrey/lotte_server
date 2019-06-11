package model

// DbChannel ...
type DbChannel chan DbResult

// DbResult DB 응답 타입 정의
type DbResult struct {
	Data interface{}
	Err  *AppError
}

// AppError Error정보
type AppError struct {
	ErrorCode  int32  `json:"errorCode"`
	ErrorMsg   string `json:"errorMsg"`
	Message    string `json:"message"`
	StatusCode int32  `json:"statusCode"`
	CreatedAt  int64  `json:"createdAt"`
}

// JSONResponse ...
type JSONResponse struct {
	Error     *AppError    `json:"error"`
	Result    ResponseData `json:"result"`
	Timestamp int64        `json:"timestamp"`
	Success   bool         `json:"success"`
}

// ResponseData ...
type ResponseData struct {
	Data interface{} `json:"data"`
}