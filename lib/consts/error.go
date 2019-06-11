package consts

// 에러 정의
const (
	ErrorBadRequestMsg  = "필수 항목 값이 없습니다."
	ErrorBadRequestCode = 1000

	ErrorCreateUserMsg  = "회원 정보 생성이 실패했습니다."
	ErrorCreateUserCode = 1001

	ErrorFindUserMsg  = "등록된 회원 정보가 없습니다."
	ErrorFindUserCode = 1002

	ErrorValidatePasswordMsg  = "비밀번호가 틀립니다."
	ErrorValidatePasswordCode = 1003
)
