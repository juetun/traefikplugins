package logic

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/juetun/traefikplugins/logic/sign"
	"github.com/juetun/traefikplugins/pkg"
)

// Config the plugin configuration.
type (
	Config struct {
		NotNeedLogin   bool              `json:"not_need_login"` // 不需要登录验证 默认需要登录验证
		NotNeedSign    bool              `json:"not_need_sign"`  // 不需要签名验证 默认需要签名验证
		PermitValidate string            `json:"permit_validate"`
		AppEnv         string            `json:"app_env,omitempty"` // 运行环境
		RouterType     string            `json:"router_type,omitempty"`
		Headers        map[string]string `json:"headers,omitempty"`
	}
	RouteTypeBaseLogic struct {
		Response     http.ResponseWriter `json:"-"`
		Request      *http.Request       `json:"-"`
		Next         http.Handler        `json:"-"`
		Ctx          context.Context     `json:"-"`
		PluginConfig *Config             `json:"plugin_config"`
		Name         string              `json:"name"`
		UserHid      string              `json:"user_hid"`  // 当前登录用户ID
		NeedSign     bool                `json:"need_sign"` // 默认需要验证
		UriParam     UriParam            `json:"uri"`       // 当前接口的访问路径
	}
	UriParam struct {
		Method   string `json:"method"`
		AppName  string `json:"app_name"`
		PathType string `json:"path_type"`
		Uri      string `json:"uri"`
	}
)

// FlagNeedLogin TODO 判断当前是否需要登录
func (r *RouteTypeBaseLogic) FlagNeedLogin() (needLogin bool, errCode int) {
	if r.PluginConfig.NotNeedLogin { // 如果不需要判断登录
		return
	}

	return
}

// FlagHavePermit 判断是否有权限使用接口
func (r *RouteTypeBaseLogic) FlagHavePermit() (errCode int) {
	if r.PluginConfig.PermitValidate == "" {
		return
	}
	return
}

func (r *RouteTypeBaseLogic) ParseUriParam() (errCode int, errMessage string) {
	r.UriParam = UriParam{
		Method: r.Request.Method,
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
	return
}

// CommonLogic 公共的逻辑模块
func (r *RouteTypeBaseLogic) CommonLogic() (exit bool) {
	var (
		err     error
		errCode int
		errMsg  string
		ok      bool
	)
	if errCode, errMsg = r.ParseUriParam(); errCode != 0 {
		if errMsg == "" {
			errMsg = pkg.MapGatewayError[errCode]
		}
		exit = true
		http.Error(r.Response, errMsg, http.StatusOK)
		return
	}
	if errCode, errMsg = r.SignValidate(); errCode != 0 {
		if errMsg == "" {
			errMsg = pkg.MapGatewayError[errCode]
		}
		exit = true
		http.Error(r.Response, errMsg, http.StatusOK)
		return
	}
	// 如果不需要登录
	if ok, errCode = r.FlagNeedLogin(); errCode != 0 {
		exit = true
		http.Error(r.Response, pkg.MapGatewayError[errCode], http.StatusOK)
		return
	} else if !ok { // 如果不需要登录
		if r.PluginConfig.PermitValidate != "" {
			errCode = pkg.GatewayErrorCodePermitConfigError
			http.Error(r.Response, fmt.Sprintf(pkg.MapGatewayError[errCode], r.Name), http.StatusOK)
			return
		}
		return
	}

	// 判断是否登录
	if errCode = r.FlagLogin(); errCode != 0 {
		exit = true
		http.Error(r.Response, err.Error(), http.StatusOK)
		return
	}
	// 判断是否有权限使用接口
	if errCode = r.FlagHavePermit(); errCode != 0 {
		exit = true
		http.Error(r.Response, err.Error(), http.StatusOK)
		return
	}
	return
}

// FlagLogin TODO 判断当前是否登录
func (r *RouteTypeBaseLogic) FlagLogin() (errCode int) {
	if r.PluginConfig.NotNeedLogin { // 如果不需要判断登录
		return
	}
	return
}

func (r *RouteTypeBaseLogic) SignValidate() (errCode int, errMsg string) {
	if r.PluginConfig.NotNeedSign { // 如果不需要判断签名验证
		return
	}
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
	res = &RouteTypeBaseLogic{
		NeedSign: true,
	}
	for _, option := range options {
		option(res)
	}
	return
}
