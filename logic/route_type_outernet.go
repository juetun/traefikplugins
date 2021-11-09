package logic

type (
	RouteTypeOuternetLogic struct {
		*RouteTypeBaseLogic
	}
	OptionsOuternetHandler func(rbc *RouteTypeOuternetLogic)
)

func (r *RouteTypeOuternetLogic) Run() {

	r.NextExecute()
}

// FlagHavePermit 判断是否有权限使用接口
func (r *RouteTypeOuternetLogic) FlagHavePermit() (errCode int) {

	return
}
func OptionsOuternetHandlerBase(rTBl *RouteTypeBaseLogic) OptionsOuternetHandler {
	return func(rbc *RouteTypeOuternetLogic) {
		rbc.RouteTypeBaseLogic = rTBl
	}
}
func NewRouteTypeOuternetLogic(options ...OptionsOuternetHandler) (res *RouteTypeOuternetLogic) {
	res = &RouteTypeOuternetLogic{}
	for _, option := range options {
		option(res)
	}

	return
}
