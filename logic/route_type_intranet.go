package logic


type (
	RouteTypeIntranetLogic struct {
		*RouteTypeBaseLogic
	}
	OptionsIntranetHandler func(rbc *RouteTypeIntranetLogic)
)

func (r *RouteTypeIntranetLogic) Run() {
	r.NextExecute()
}



func OptionsIntranetHandlerBase(rTBl *RouteTypeBaseLogic) OptionsIntranetHandler {
	return func(rbc *RouteTypeIntranetLogic) {
		rbc.RouteTypeBaseLogic = rTBl
	}
}

func NewRouteTypeIntranetLogic(options ...OptionsIntranetHandler) (res *RouteTypeIntranetLogic) {
	res = &RouteTypeIntranetLogic{}
	for _, option := range options {
		option(res)
	}

	return
}
