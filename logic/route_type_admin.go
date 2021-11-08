package logic

type (
	RouteTypeAdminLogic struct {
		*RouteTypeBaseLogic
	}
	OptionsAdminHandler func(rbc *RouteTypeAdminLogic)
)

func (r *RouteTypeAdminLogic) Run() {
 	r.NextExecute()
}



func OptionsAdminHandlerBase(rTBl *RouteTypeBaseLogic) OptionsAdminHandler {
	return func(rbc *RouteTypeAdminLogic) {
		rbc.RouteTypeBaseLogic = rTBl
	}
}
func NewRouteTypeAdminLogic(options ...OptionsAdminHandler) (res *RouteTypeAdminLogic) {
	res = &RouteTypeAdminLogic{}
	for _, option := range options {
		option(res)
	}

	return
}
