package pkg

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
		UserHid  string `json:"user_hid"` // 当前登录用户ID
		Method   string `json:"method"`
		AppName  string `json:"app_name"`
		PathType string `json:"path_type"`
		Uri      string `json:"uri"`
	}
)
