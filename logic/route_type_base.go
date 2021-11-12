package logic

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/juetun/traefikplugins/logic/sign"
	"github.com/juetun/traefikplugins/pkg"
)

var (
	GrpcGet            pkg.PermitGet
	HttpGet            pkg.PermitGet
	ConfigRouterPermit pkg.RouterConfig // 当前系统支持的路由权限
)

// Config the plugin configuration.
type (
	Config struct {
		PermitSupportGrpc bool              `json:"permit_support_grpc"` // 获取接口权限是否支持grpc
		PermitValidate    string            `json:"permit_validate"`     // 是否需要签名验证
		AppEnv            string            `json:"app_env,omitempty"`   // 运行环境
		RouterType        string            `json:"router_type,omitempty"`
		Headers           map[string]string `json:"headers,omitempty"`
	}
	RouteTypeBaseLogic struct {
		Response     http.ResponseWriter `json:"-"`
		Request      *http.Request       `json:"-"`
		Next         http.Handler        `json:"-"`
		Ctx          context.Context     `json:"-"`
		PluginConfig *Config             `json:"plugin_config"`
		Name         string              `json:"name"`
		UriParam     pkg.UriParam        `json:"uri"` // 当前接口的访问路径
	}
)

func (r *RouteTypeBaseLogic) ParseUriParam() (errCode int, errMessage string) {
	r.UriParam = pkg.UriParam{
		Method:  r.Request.Method,
		UserHid: r.Request.Header.Get(pkg.HttpUserHid),
	}
	const divString = "/"
	var pathSlice = make([]string, 0, 10)
	path := strings.Split(r.Request.RequestURI, divString)
	l := len(path)
	pathSlice = path[0:]
	switch l {
	case 0:
		pathSlice[0] = ""
		pathSlice[1] = ""
		pathSlice[2] = ""
	case 1:
		pathSlice[1] = ""
		pathSlice[2] = ""
	case 2:
		pathSlice[2] = ""
	}
	r.UriParam.AppName = pathSlice[0]
	r.UriParam.PathType = pathSlice[1]
	pathSlice = pathSlice[2:]
	r.UriParam.Uri = strings.Join(pathSlice, divString)
	if r.UriParam.AppName == "" || r.UriParam.PathType == "" || r.UriParam.Uri == "" {
		errCode = pkg.GateWayPathError
		errMessage = fmt.Sprintf(pkg.MapGatewayError[errCode], r.Request.RequestURI)
		return
	}
	return
}

func (r *RouteTypeBaseLogic) LoadUrlConfig() (errCode int, errMsg string) {

	if r.UriParam.AppName !=  pkg.RouteTypeGateway {
		return
	}
	if r.UriParam.Method != http.MethodHead {
		errCode = pkg.GateWayLoadConfigError
		errMsg = fmt.Sprintf(pkg.MapGatewayError[errCode], "METHOD")
		return
	}
	if r.UriParam.Uri != "" {

	}
	var newConfigRouter pkg.RouterConfig
	var lock sync.RWMutex
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	ConfigRouterPermit = newConfigRouter
	return
}

// CommonLogic 公共的逻辑模块
func (r *RouteTypeBaseLogic) CommonLogic() (exit bool) {
	switch r.PluginConfig.PermitValidate {
	case "": // 如果都不需要验证操作，则直接跳过
		exit = true
		http.Error(r.Response, fmt.Sprintf(pkg.MapGatewayError[pkg.GatewayErrorCodePermitConfigError], "permit_validate is null"), http.StatusInternalServerError)
		return
	}

	var (
		errCode int
		errMsg  string
	)

	// 拆分Path路径
	if errCode, errMsg = r.ParseUriParam(); errCode != 0 {
		if errMsg == "" {
			errMsg = pkg.MapGatewayError[errCode]
		}
		exit = true
		http.Error(r.Response, errMsg, http.StatusOK)
		return
	}

	// 判断是否不需要登录
	if errCode, errMsg = r.FlagHavePermit(); errCode != 0 {
		exit = true
		if errMsg == "" {
			errMsg = pkg.MapGatewayError[errCode]
		}
		http.Error(r.Response, pkg.MapGatewayError[errCode], http.StatusOK)
		return
	}

	return
}

// FlagHavePermit 判断是否有权限使用接口
func (r *RouteTypeBaseLogic) FlagHavePermit() (errCode int, errMsg string) {
	var res pkg.PermitGetResult
	if r.PluginConfig.PermitSupportGrpc { // 获取接口需要验证的权限
		res, errCode, errMsg = GrpcGet.Do(&r.UriParam)
	} else { // 获取接口需要验证的权限
		res, errCode, errMsg = HttpGet.Do(&r.UriParam)
	}
	if errCode != 0 {
		return
	}

	if res.NeedSign {
		// 接口签名验证判断
		errCode, errMsg = r.SignValidate()
		if errCode != 0 {
			return
		}
	}
	return
}

func (r *RouteTypeBaseLogic) SignValidate() (errCode int, errMsg string) {

	errCode, errMsg = sign.NewLogicSign(
		sign.OptionAppEnv(r.PluginConfig.AppEnv),
		sign.OptionResponse(r.Response),
		sign.OptionRequest(r.Request),
	).SignValidate()
	return
}

// NextExecute 继续往下执行
func (r *RouteTypeBaseLogic) NextExecute() {
	r.Next.ServeHTTP(r.Response, r.Request)
}

// WriteResponseHeader 向响应信息中回写header
func (r *RouteTypeBaseLogic) WriteResponseHeader(header map[string]string) {
	responseHeader := r.Response.Header()
	for key, val := range header {
		responseHeader.Add(key, val)
	}
}

// 获取接口权限
func (r *RouteTypeBaseLogic) getImportPermit() (errCode int, errMsg string) {

	return
}

type OptionsHandler func(rbc *RouteTypeBaseLogic)

func OptionsHandlerPluginCtx(ctx context.Context) OptionsHandler {
	return func(rbc *RouteTypeBaseLogic) {
		rbc.Ctx = ctx
	}
}
func OptionsHandlerPluginName(name string) OptionsHandler {
	return func(rbc *RouteTypeBaseLogic) {
		rbc.Name = name
	}
}
func OptionsHandlerPluginConfig(PluginConfig *Config) OptionsHandler {
	return func(rbc *RouteTypeBaseLogic) {
		rbc.PluginConfig = PluginConfig
	}
}
func OptionsHandlerNext(Next http.Handler) OptionsHandler {
	return func(rbc *RouteTypeBaseLogic) {
		rbc.Next = Next
	}
}
func OptionsHandlerRequest(request *http.Request) OptionsHandler {
	return func(rbc *RouteTypeBaseLogic) {
		rbc.Request = request
	}
}
func OptionsHandlerResponse(Response http.ResponseWriter) OptionsHandler {
	return func(rbc *RouteTypeBaseLogic) {
		rbc.Response = Response
	}
}

func NewRouteTypeBaseLogic(options ...OptionsHandler) (res *RouteTypeBaseLogic) {
	res = &RouteTypeBaseLogic{}
	for _, option := range options {
		option(res)
	}
	return
}
