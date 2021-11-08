// Package traefikplugins
// Package plugindemo a demo plugin.
package traefikplugins

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/juetun/traefikplugins/logic"
	"github.com/juetun/traefikplugins/pkg"
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
	switch r.PluginConfig.RouterType {
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
