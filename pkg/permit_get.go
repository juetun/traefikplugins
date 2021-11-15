package pkg

import (
	"regexp"
)

// PermitGet 获取权限接口操作
type (
	PermitGet interface {
		Do(uriParam *UriParam) (permitResult PermitGetResult, errCode int, errMsg string)
	}
	PermitGetResult struct {
		NeedLogin  bool `json:"need_login"`
		NeedSign   bool `json:"need_sign"`
		HavePermit bool `json:"have_permit"`
	}

	UriParam struct {
		UserHid   string `json:"user_hid"` // 当前登录用户ID
		Method    string `json:"method"`
		AppName   string `json:"app_name"`
		PathType  string `json:"path_type"`
		RegexpUri string `json:"regexp_uri"`
		Uri       string `json:"uri"`
	}

	RouterConfig struct {
		RouterNotNeedSign  map[string]*RouterNotNeedItem `json:"not_sign"`  // 不需要签名验证的路由权限
		RouterNotNeedLogin map[string]*RouterNotNeedItem `json:"not_login"` // 不需要登录的路由权限
	}
	RouterNotNeedItem struct {
		GeneralPath map[string]ItemGateway `json:"general,omitempty"` // 普通路径
		RegexpPath  []ItemGateway          `json:"regexp,omitempty"`  // 按照正则匹配的路径
	}
	ItemGateway struct {
		Uri     string   `json:"url,omitempty"`
		Methods []string `json:"method,omitempty"`
	}
)

// MathValidateType 获取当前路径的校验规则
func (r *RouterConfig) MathValidateType(uriParam *UriParam) (notNeedSign, notNeedLogin bool, err error) {

	if notNeedSign, err = r.notNeedSignValidate(uriParam); err != nil {
		return
	}
	if notNeedLogin, err = r.notNeedLoginValidate(uriParam); err != nil {
		return
	}

	return
}

func (r *RouterConfig) notNeedSignValidate(uriParam *UriParam) (notNeedSign bool, err error) {
	var (
		dt    *RouterNotNeedItem
		ok    bool
		match bool
	)
	if dt, ok = r.RouterNotNeedSign[uriParam.AppName]; !ok {
		return
	}
	if _, ok = dt.GeneralPath[uriParam.RegexpUri]; ok {
		notNeedSign = true
		return
	}

	for _, item := range dt.RegexpPath {
		if match, err = r.routePathMath(&item, uriParam.RegexpUri); err != nil {
			return
		}
		if match {
			notNeedSign = true
			return
		}
	}
	return
}
func (r *RouterConfig) notNeedLoginValidate(uriParam *UriParam) (notNeedLogin bool, err error) {

	var (
		dt    *RouterNotNeedItem
		ok    bool
		match bool
	)

	if dt, ok = r.RouterNotNeedLogin[uriParam.AppName]; !ok {
		return
	}

	if _, ok = dt.GeneralPath[uriParam.RegexpUri]; ok {
		notNeedLogin = true
		return
	}

	for _, item := range dt.RegexpPath {
		if match, err = r.routePathMath(&item, uriParam.RegexpUri); err != nil {
			return
		}
		if match {
			notNeedLogin = true
			return
		}
	}

	return
}
func (r *RouterConfig) routePathMath(reg *ItemGateway, path string) (matched bool, err error) {
	matched, err = regexp.Match(reg.Uri, []byte(path))
	return
}
