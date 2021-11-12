// Package traefikplugins
// Package plugindemo a demo plugin.
package traefikplugins

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/juetun/traefikplugins/logic"
	"github.com/juetun/traefikplugins/pkg"
	"github.com/juetun/traefikplugins/pkg/permit_get"
)

// TRaeFikJueTun 网站插件
type TRaeFikJueTun struct {
	Next         http.Handler
	PluginConfig *logic.Config
	Name         string
	Ctx          context.Context
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *logic.Config {
	return &logic.Config{
		Headers: make(map[string]string),
	}
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *logic.Config, name string) (httpHandler http.Handler, err error) {
	httpHandler = &TRaeFikJueTun{
		Ctx:          ctx,
		PluginConfig: config,
		Next:         next,
		Name:         name,
	}

	// 初始化获取权限操作对象
	logic.GrpcGet = permit_get.NewGrpcGet()
	logic.HttpGet = permit_get.NewHttpGet()

	return
}
func (r *TRaeFikJueTun) ToJsonString() (res string, err error) {
	var configBt []byte
	if configBt, err = json.Marshal(r.PluginConfig); err != nil {
		return
	}
	res = string(configBt)
	return
}

func (r *TRaeFikJueTun) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	for s, s2 := range r.PluginConfig.Headers {
		request.Header.Set(s, s2)
	}
	switch r.PluginConfig.RouterType {
	case pkg.RouteTypeGateway: // 此路径只用与更新路由规则调用
		logicOp := logic.NewRouteTypeAdminLogic(logic.OptionsAdminHandlerBase(r.getBaseLogic(response, request)))
		if errCode, errMsg := logicOp.LoadUrlConfig(); errCode != 0 {
			http.Error(response, errMsg, http.StatusInternalServerError)
			return
		}
		type Result struct {
			Code int         `json:"code"`
			Data interface{} `json:"data"`
			Msg  string      `json:"message"`
		}
		var Res = Result{
			Msg: "路由规则更新成功",
		}
		var bt []byte
		bt, _ = json.Marshal(Res)
		http.Error(response, string(bt), http.StatusOK)
	case pkg.RouteTypeAdmin:
		logicOp := logic.NewRouteTypeAdminLogic(logic.OptionsAdminHandlerBase(r.getBaseLogic(response, request)))
		if exit := logicOp.CommonLogic(); exit {
			return
		}
		logicOp.Run()
	case pkg.RouteTypeIntranet:
		logicOp := logic.NewRouteTypeIntranetLogic(logic.OptionsIntranetHandlerBase(r.getBaseLogic(response, request)))
		if exit := logicOp.CommonLogic(); exit {
			return
		}
		logicOp.Run()
	case pkg.RouteTypeOuternet:
		logicOp := logic.NewRouteTypeOuternetLogic(logic.OptionsOuternetHandlerBase(r.getBaseLogic(response, request)))
		if exit := logicOp.CommonLogic(); exit {
			return
		}
		logicOp.Run()
	case pkg.RouteTypePage:
		logicOp := logic.NewRouteTypePageLogic(logic.OptionsPageHandlerBase(r.getBaseLogic(response, request)))
		if exit := logicOp.CommonLogic(); exit {
			return
		}
		logicOp.Run()
	default:
		r.Next.ServeHTTP(response, request)
	}
}

func (r *TRaeFikJueTun) getBaseLogic(response http.ResponseWriter, request *http.Request) (base *logic.RouteTypeBaseLogic) {
	base = logic.NewRouteTypeBaseLogic(
		logic.OptionsHandlerPluginCtx(r.Ctx),
		logic.OptionsHandlerPluginName(r.Name),
		logic.OptionsHandlerPluginConfig(r.PluginConfig),
		logic.OptionsHandlerNext(r.Next),
		logic.OptionsHandlerRequest(request),
		logic.OptionsHandlerResponse(response),
	)
	return
}
