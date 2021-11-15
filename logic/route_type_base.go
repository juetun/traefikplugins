package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/juetun/traefikplugins/logic/sign"
	"github.com/juetun/traefikplugins/pkg"
)

var (
	GrpcGet            pkg.PermitGet
	HttpGet            pkg.PermitGet
	ConfigRouterPermit *pkg.RouterConfig // 当前系统支持的路由权限
)

// Config the plugin configuration.
type (
	Config struct {
		TraefikConfigPluginName string     `json:"traefik_config_plugin_name"`
		AppEnv                  string     `json:"appenv,omitempty"` // 运行环境
		RouterType              string     `json:"routertype,omitempty"`
		PermitSupperFlag        uint8      `json:"supperflag,omitempty"` // 是否客服后台权限验证
		PathConfig              PathConfig `json:"pathconfig,omitempty"`
	}
	PathConfig struct {
		Host    string `json:"host,omitempty"` // 获取不需要签名验证和登录验证的接口地址调用接口来源host
		AppName string `json:"app,omitempty"`  // 获取不需要签名验证和登录验证的接口地址调用接口应用路径
	}
	RouteTypeBaseLogic struct {
		Response     http.ResponseWriter `json:"-"`
		Request      *http.Request       `json:"-"`
		Next         http.Handler        `json:"-"`
		Ctx          context.Context     `json:"-"`
		PluginConfig *Config             `json:"plugin_config"`
		UriParam     *pkg.UriParam       `json:"uri"` // 当前接口的访问路径
	}
)

func (r *RouteTypeBaseLogic) LoadUrlConfig() (errCode int, errMsg string) {

	if r.UriParam.AppName != pkg.RouteTypeGateway {
		return
	}

	if r.UriParam.Method != http.MethodHead {
		errCode = pkg.GateWayLoadConfigError
		errMsg = fmt.Sprintf(pkg.MapGatewayError[errCode], "METHOD")
		return
	}

	if r.UriParam.Uri != "/load_config" {
		errCode = pkg.GateWayLoadConfigError
		errMsg = fmt.Sprintf(pkg.MapGatewayError[errCode], "当前不支持你访问的接口路径")
		return
	}
	errCode, errMsg = r.RefreshConfigRouterPermit()
	return
}
func (r *RouteTypeBaseLogic) RefreshConfigRouterPermit(appNames ...string) (errCode int, errMsg string) {
	var newConfigRouter *pkg.RouterConfig
	if newConfigRouter, errCode, errMsg = r.GetUrlConfigFromDashboardAdmin(appNames...); errCode != 0 {
		return
	}
	// 更新路由配置数据时加锁 防止数据串改
	var lock sync.RWMutex
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()

	// 组织不需要签名的接口列表
	if newConfigRouter.RouterNotNeedSign != nil {
		for appName, item := range newConfigRouter.RouterNotNeedSign {
			if ConfigRouterPermit.RouterNotNeedSign == nil {
				ConfigRouterPermit.RouterNotNeedSign = map[string]*pkg.RouterNotNeedItem{}
			}
			ConfigRouterPermit.RouterNotNeedSign[appName] = item
		}
	}

	// 组织不需要登录的接口列表
	if newConfigRouter.RouterNotNeedLogin != nil {
		for appName, item := range newConfigRouter.RouterNotNeedLogin {
			if ConfigRouterPermit.RouterNotNeedLogin == nil {
				ConfigRouterPermit.RouterNotNeedLogin = map[string]*pkg.RouterNotNeedItem{}
			}
			ConfigRouterPermit.RouterNotNeedLogin[appName] = item
		}
	}
	return
}
func (r *RouteTypeBaseLogic) GetUrlConfigFromDashboardAdmin(appName ...string) (res *pkg.RouterConfig, errCode int, errMsg string) {
	type MyJsonName struct {
		Code    int64             `json:"code"`
		Data    *pkg.RouterConfig `json:"data"`
		Message string            `json:"message"`
	}

	var (
		err error
		Res MyJsonName
	)

	if err = r.HttpGetUrlConfig(&Res, appName...); err != nil {
		return
	}
	res = Res.Data
	return
}
func (r *RouteTypeBaseLogic) getAddr(apps ...string) (reqAddress string) {
	appName := strings.Join(apps, ",")
	if appName == "" {
		reqAddress = fmt.Sprintf("%s/%s/in/get_import_permit",
			r.PluginConfig.PathConfig.Host,
			r.PluginConfig.PathConfig.AppName)
	} else {
		reqAddress = fmt.Sprintf(
			"%s/%s/in/get_import_permit?app_name=%s",
			r.PluginConfig.PathConfig.Host,
			r.PluginConfig.PathConfig.AppName,
			appName,
		)
	}
	return
}
func (r *RouteTypeBaseLogic) HttpGetUrlConfig(data interface{}, apps ...string) (err error) {

	var reqAddress = r.getAddr(apps...)
	var client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, netW, addr string) (conn net.Conn, e error) {
				conn, e = net.DialTimeout(netW, addr, time.Second*3)
				if e != nil {
					return
				}
				return
			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 2,
		},
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, reqAddress, nil)
	var resp *http.Response
	var body []byte
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(body, data); err != nil {
		return
	}
	return
}

// CommonLogic 公共的逻辑模块
func (r *RouteTypeBaseLogic) CommonLogic() (exit bool) {

	var (
		errCode int
		errMsg  string
	)

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

func (r *RouteTypeBaseLogic) notNeedSignLogic() (notNeedSign bool) {

	var (
		ok   bool
		dt   *pkg.RouterNotNeedItem
		item pkg.ItemGateway
	)

	if dt, ok = ConfigRouterPermit.RouterNotNeedSign[r.UriParam.AppName]; !ok {
		return
	}

	if item, ok = dt.GeneralPath[r.UriParam.Uri]; ok {
		if _, ok = item.Methods[r.UriParam.Method]; ok {
			notNeedSign = true
			return
		}
		return
	}

	for _, item = range dt.RegexpPath {
		if ok, _ = r.RoutePathMath(item.Uri, r.UriParam.Uri); !ok {
			continue
		}
		if _, ok = item.Methods[r.UriParam.Method]; ok {
			notNeedSign = true
			return
		}
	}
	return
}

func (r *RouteTypeBaseLogic) notNeedLoginLogic() (notNeedLogin bool) {
	var (
		ok   bool
		dt   *pkg.RouterNotNeedItem
		item pkg.ItemGateway
	)

	if dt, ok = ConfigRouterPermit.RouterNotNeedLogin[r.UriParam.AppName]; !ok {
		return
	}

	if item, ok = dt.GeneralPath[r.UriParam.Uri]; ok {
		if _, ok = item.Methods[r.UriParam.Method]; ok {
			notNeedLogin = true
			return
		}
		return
	}
	for _, item = range dt.RegexpPath {
		if ok, _ = r.RoutePathMath(item.Uri, r.UriParam.Uri); !ok {
			continue
		}
		if _, ok = item.Methods[r.UriParam.Method]; ok {
			notNeedLogin = true
			return
		}
	}

	return
}
func (r *RouteTypeBaseLogic) RoutePathMath(regexpString, path string) (matched bool, err error) {
	matched, err = regexp.Match(regexpString, []byte(path))
	return
}

// FlagHavePermit 判断是否有权限使用接口
func (r *RouteTypeBaseLogic) FlagHavePermit() (errCode int, errMsg string) {

	// 接口签名验证判断
	if notNeedSign := r.notNeedSignLogic(); !notNeedSign {
		if errCode, errMsg = r.SignValidate(); errCode != 0 {
			return
		}
	}

	// 接口登录验证判断
	if notNeedLogin := r.notNeedLoginLogic(); !notNeedLogin {
		if errCode, errMsg = r.SignValidate(); errCode != 0 {
			return
		}
	}

	if r.PluginConfig.PermitSupperFlag > 0 { // 获取接口需要验证的权限
		if _, errCode, errMsg = HttpGet.Do(r.UriParam); errCode != 0 {
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

// func OptionsHandlerPluginName(name string) OptionsHandler {
// 	return func(rbc *RouteTypeBaseLogic) {
// 		rbc.Name = name
// 	}
// }

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
func OptionsHandlerUrlParam(uriParam *pkg.UriParam) OptionsHandler {
	return func(rbc *RouteTypeBaseLogic) {
		rbc.UriParam = uriParam
	}
}

func NewRouteTypeBaseLogic(options ...OptionsHandler) (res *RouteTypeBaseLogic) {
	res = &RouteTypeBaseLogic{}
	for _, option := range options {
		option(res)
	}
	return
}
