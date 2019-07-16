package errno
/*
服务级别错误：1 为系统级错误；2 为普通错误，通常是由用户非法操作引起的
服务模块为两位数：一个大型系统的服务模块通常不超过两位数，如果超过，说明这个系统该拆分了
错误码为两位数：防止一个模块定制过多的错误码，后期不好维护
code = 0 说明是正确返回，code > 0 说明是错误返回
错误通常包括系统级错误码和服务级错误码
建议代码中按服务模块将错误分类
错误码均为 >= 0 的数
*/
var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	// user errors
	ErrUserNotFound      = &Errno{Code: 20102, Message: "The user was not found."}
)