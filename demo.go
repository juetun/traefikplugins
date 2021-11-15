// Package traefikplugins
// Package plugindemo a demo plugin.
package traefikplugins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/juetun/traefikplugins/logic"
	"github.com/juetun/traefikplugins/pkg"
	"github.com/juetun/traefikplugins/pkg/permit_get"
)

var mapHandlerConfig = map[string]*logic.Config{}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *logic.Config {
	return &logic.Config{
		PathConfig: logic.PathConfig{},
	}
}

// TRaeFikJueTun 网站插件
type TRaeFikJueTun struct {
	Next http.Handler
	Ctx  context.Context
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *logic.Config, name string) (httpHandler http.Handler, err error) {

	tRaeFikJueTun := &TRaeFikJueTun{
		Ctx:  ctx,
		Next: next,
	}

	// 初始化获取接口不需要登录，和不需要签名验证的接口列表
	if err = tRaeFikJueTun.PreloadImportConfig(config); err != nil {
		return
	}

	config.TraefikConfigPluginName = name
	mapHandlerConfig[config.RouterType] = config

	// 初始化获取权限操作对象
	logic.GrpcGet = permit_get.NewGrpcGet()
	logic.HttpGet = permit_get.NewHttpGet()
	return tRaeFikJueTun, err
}

// PreloadImportConfig 预状态不需要签名验证 和不需要登录的接口列表
func (r *TRaeFikJueTun) PreloadImportConfig(config *logic.Config) (err error) {
	if logic.ConfigRouterPermit != nil {
		return
	}
	routeTypeBaseLogic := logic.RouteTypeBaseLogic{PluginConfig: config}
	errCode, errMsg := routeTypeBaseLogic.RefreshConfigRouterPermit()
	if errCode != 0 {
		err = fmt.Errorf(errMsg)
	}
	return
}

func (r *TRaeFikJueTun) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	var (
		urlParam *pkg.UriParam
		errCode  int
		errMsg   string
	)
	if urlParam, errCode, errMsg = r.ParseUriParam(request); errCode != 0 {
		http.Error(response, errMsg, http.StatusInternalServerError)
		return
	}
	switch urlParam.PathType {
	case pkg.RouteTypeAdmin:
		logicOp := logic.NewRouteTypeAdminLogic(logic.OptionsAdminHandlerBase(r.getBaseLogic(response, request, urlParam, mapHandlerConfig[pkg.RouteTypeAdmin])))
		var exit bool
		if exit = r.loadConfig(logicOp, response); exit {
			return
		}
		if exit = logicOp.CommonLogic(); exit {
			return
		}
		logicOp.Run()
	case pkg.RouteTypeIntranet:
		logicOp := logic.NewRouteTypeIntranetLogic(logic.OptionsIntranetHandlerBase(r.getBaseLogic(response, request, urlParam, mapHandlerConfig[pkg.RouteTypeIntranet])))
		if exit := logicOp.CommonLogic(); exit {
			return
		}
		logicOp.Run()
	case pkg.RouteTypeOuternet:
		logicOp := logic.NewRouteTypeOuternetLogic(logic.OptionsOuternetHandlerBase(r.getBaseLogic(response, request, urlParam, mapHandlerConfig[pkg.RouteTypeOuternet])))
		if exit := logicOp.CommonLogic(); exit {
			return
		}
		logicOp.Run()
	case pkg.RouteTypePage:
		logicOp := logic.NewRouteTypePageLogic(logic.OptionsPageHandlerBase(r.getBaseLogic(response, request, urlParam, mapHandlerConfig[pkg.RouteTypePage])))
		if exit := logicOp.CommonLogic(); exit {
			return
		}
		logicOp.Run()
	default:
		r.Next.ServeHTTP(response, request)
	}
}

func (r *TRaeFikJueTun) loadConfig(logicOp *logic.RouteTypeAdminLogic, response http.ResponseWriter) (exit bool) {

	if logicOp.UriParam.AppName != pkg.RouteTypeGateway {
		return
	}
	type Result struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
		Msg  string      `json:"message"`
	}
	var Res = Result{Msg: "路由规则更新成功"}
	defer func() {
		if exit {
			var bt []byte
			bt, _ = json.Marshal(Res)
			http.Error(response, string(bt), http.StatusOK)
		}
	}()
	if errCode, errMsg := logicOp.LoadUrlConfig(); errCode != 0 {
		exit = true
		Res.Msg = errMsg
		return
	}
	exit = true
	return
}

func (r *TRaeFikJueTun) getBaseLogic(response http.ResponseWriter, request *http.Request, urlParam *pkg.UriParam, config *logic.Config) (base *logic.RouteTypeBaseLogic) {
	base = logic.NewRouteTypeBaseLogic(
		logic.OptionsHandlerPluginCtx(r.Ctx),
		// logic.OptionsHandlerPluginName(r.Name),
		logic.OptionsHandlerPluginConfig(config),
		logic.OptionsHandlerNext(r.Next),
		logic.OptionsHandlerRequest(request),
		logic.OptionsHandlerResponse(response),
		logic.OptionsHandlerUrlParam(urlParam),
	)
	return
}

// ParseUriParam 解析path路径
func (r *TRaeFikJueTun) ParseUriParam(request *http.Request) (uriParam *pkg.UriParam, errCode int, errMessage string) {
	uriParam = &pkg.UriParam{
		Method:  request.Method,
		UserHid: request.Header.Get(pkg.HttpUserHid),
	}
	const divString = "/"
	var pathSlice = make([]string, 0, 10)

	path := strings.Split(request.RequestURI, divString)
	l := len(path)
	pathSlice = path[0:]
	switch l {
	case 0:
		pathSlice[0] = ""
		pathSlice[1] = ""
		pathSlice[2] = ""
		pathSlice[3] = ""
	case 1:
		pathSlice[1] = ""
		pathSlice[2] = ""
		pathSlice[3] = ""
	case 2:
		pathSlice[2] = ""
		pathSlice[3] = ""
	case 3:
		pathSlice[3] = ""
	}
	uriParam.AppName = pathSlice[1]
	uriParam.PathType = pathSlice[2]
	pathSlice = pathSlice[3:]
	uriParam.Uri = strings.Join(pathSlice, divString)
	if uriParam.AppName == "" || uriParam.PathType == "" || uriParam.Uri == "" {
		errCode = pkg.GateWayPathError
		errMessage = fmt.Sprintf(pkg.MapGatewayError[errCode], request.RequestURI)
		return
	}
	return
}
