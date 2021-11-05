package logic

import (
	"net/http"

	"github.com/juetun/traefikplugins/pkg"
)

type (
	RouteTypeOuternetLogic struct {
		*RouteTypeBaseLogic
	}
	OptionsOuternetHandler func(rbc *RouteTypeOuternetLogic)
)

func (r *RouteTypeOuternetLogic) Run() {
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
