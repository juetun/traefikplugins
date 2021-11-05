package logic

import (
	"net/http"

	"github.com/juetun/traefikplugins/pkg"
)

type (
	RouteTypePageLogic struct {
		*RouteTypeBaseLogic
	}
	OptionsPageHandler func(rbc *RouteTypePageLogic)
)

func (r *RouteTypePageLogic) Run() {
	var err error
	var errCode int
	var ok bool

	// 网页访问没有签名验证
	// if errCode = r.SignValidate(); errCode != 0 {
	// 	http.Error(r.Response, pkg.MapGatewayError[errCode], http.StatusOK)
	// 	return
	// }

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
