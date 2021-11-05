package logic

import (
	"net/http"

	"github.com/juetun/traefikplugins/pkg"
)

type (
	RouteTypeIntranetLogic struct {
		*RouteTypeBaseLogic
	}
	OptionsIntranetHandler func(rbc *RouteTypeIntranetLogic)
)

func (r *RouteTypeIntranetLogic) Run() {
	var (
		err     error
		errCode int
		errMsg  string
		ok      bool
	)

	if errCode, errMsg = r.SignValidate(); errCode != 0 {
		if errMsg == "" {
			errMsg = pkg.MapGatewayError[errCode]
		}
		http.Error(r.Response, errMsg, http.StatusOK)
		return
	}
	// 如果不需要登录
	if ok, errCode = r.FlagNeedLogin(); errCode != 0 {
		http.Error(r.Response, pkg.MapGatewayError[errCode], http.StatusOK)
		return
	} else if !ok { // 如果不需要登录
		r.NextExecute()
		return
	}

	// 判断是否登录
	if errCode = r.FlagLogin(); errCode != 0 {
		http.Error(r.Response, err.Error(), http.StatusOK)
		return
	}
	// 判断是否有权限使用接口
	if errCode = r.FlagHavePermit(); errCode != 0 {
		http.Error(r.Response, err.Error(), http.StatusOK)
		return
	}
	r.NextExecute()
}

// FlagHavePermit 判断是否有权限使用接口
func (r *RouteTypeIntranetLogic) FlagHavePermit() (errCode int) {

	return
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
