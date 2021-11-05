package logic

import (
	"context"
	"net/http"

	"github.com/juetun/traefikplugins/logic/sign"
)

// Config the plugin configuration.
type (
	Config struct {
		AppEnv     string            `json:"app_env,omitempty"` // 运行环境
		RouterType string            `json:"router_type,omitempty"`
		Headers    map[string]string `json:"headers,omitempty"`
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
	}
)

// FlagNeedLogin TODO 判断当前是否需要登录
func (r *RouteTypeBaseLogic) FlagNeedLogin() (needLogin bool, errCode int) {

	return
}

// FlagLogin TODO 判断当前是否登录
func (r *RouteTypeBaseLogic) FlagLogin() (errCode int) {

	return
}

func (r *RouteTypeAdminLogic) SignValidate() (errCode int, errMsg string) {
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
