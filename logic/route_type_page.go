package logic

type (
	RouteTypePageLogic struct {
		*RouteTypeBaseLogic
	}
	OptionsPageHandler func(rbc *RouteTypePageLogic)
)

func (r *RouteTypePageLogic) Run() {

	r.NextExecute()
}

// FlagHavePermit 判断是否有权限使用接口
func (r *RouteTypePageLogic) FlagHavePermit() (errCode int) {

	return
}
func OptionsPageHandlerBase(rTBl *RouteTypeBaseLogic) OptionsPageHandler {
	return func(rbc *RouteTypePageLogic) {
		rbc.RouteTypeBaseLogic = rTBl
	}
}
func NewRouteTypePageLogic(options ...OptionsPageHandler) (res *RouteTypePageLogic) {
	res = &RouteTypePageLogic{}
	for _, option := range options {
		option(res)
	}

	return
}
