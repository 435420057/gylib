package respcode

type RespCode struct {
	Code string
	Info string
}

var (
	RC_GENERAL_SUCC *RespCode = &RespCode{"RC00000", "Request completed"}
	RC_GENERAL_APP_ERR *RespCode = &RespCode{"RC80000", "Application error"}
	RC_GENERAL_SYS_ERR *RespCode = &RespCode{"RC90000", "System error"}
)