package logic

import (
	"net/http"

	"github.com/juetun/traefikplugins/pkg"
)

type (
	RouteTypeAdminLogic struct {
		*RouteTypeBaseLogic
	}
	OptionsAdminHandler func(rbc *RouteTypeAdminLogic)
)

func (r *RouteTypeAdminLogic) Run() {
	var (
		errCode int
		errMsg  string
	)

	// 判断客服管理后台需要验证用户权限
	if _, errCode, errMsg = HttpGet.Do(r.UriParam); errCode != 0 {
		if errMsg == "" {
			errMsg = pkg.MapGatewayError[errCode]
		}
		http.Error(r.Response, pkg.MapGatewayError[errCode], http.StatusOK)
		return
	}
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
