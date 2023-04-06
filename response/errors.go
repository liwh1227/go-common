package response

import "fmt"

type ErrCode int32

const (
	ERR_CODE_OK                        ErrCode = iota
	ERR_CODE_INTERNAL_UNKNOWN_ERR_TYPE ErrCode = 1000
)

var errCodeName = map[ErrCode][]string{
	ERR_CODE_OK:                        {"OK", "OK"},
	ERR_CODE_INTERNAL_UNKNOWN_ERR_TYPE: {"Unknown err code type", "未知的错误类型"},
}

func (e ErrCode) String() string {
	if s, ok := errCodeName[e]; ok {
		return s[0]
	}
	return fmt.Sprintf("unknown error code %d", uint32(e))
}

func (e ErrCode) Int() int32 {
	return int32(e)
}

func (e ErrCode) CnDesc() string {
	if s, ok := errCodeName[e]; ok {
		return s[1]
	}
	return fmt.Sprintf("unkonwn error code %d", uint32(e))
}

func RegisterErrorMapper(serviceErrMap map[ErrCode][]string) {
	for k, v := range serviceErrMap {
		errCodeName[k] = v
	}
}

type Error struct {
	Code ErrCode
	Msg  string
}

func (cmErr *Error) Error() string {
	return cmErr.Msg
}
