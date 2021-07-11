package error

var (
	SUCCESS        = NewError(100000, "成功")
	NOT_FOUND      = NewError(100001, "资源不存在")
	SERVER_ERROR   = NewError(100002, "服务器错误")
	INVALID_PARAMS = NewError(100003, "请求参数错误")
)
