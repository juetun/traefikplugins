package pkg

const (
	RouteTypeGateway  = "gateway"
	RouteTypeAdmin    = "admin"
	RouteTypeIntranet = "in"
	RouteTypeOuternet = "out"
	RouteTypePage     = "page"
)
const (
	EnvDev  = "dev"  // 开发环境
	EnvTest = "test" // 测试环境
	EnvPre  = "pre"  // 预发布环境
	EnvProd = "prod" // 线上环境

	HttpUserHid = "X-User-Hid" // 页面请求时的 用户ID

)
const (
	GatewayErrorCodeSign              = iota + 10000 // 签名验证失败
	GatewayErrorCodeNotLogin                         // 未登录
	GatewayErrorCodeNotHavePermit                    // 无权限
	GatewayErrorCodePermitConfigError                // 网关配置错误
	GateWayPathError                                 // route路径错误
	GateWayLoadConfigError                           // 网关加载配置参数异常
)

var (
	MapGatewayError = map[int]string{
		GatewayErrorCodeSign:              "签名验证失败",
		GatewayErrorCodeNotLogin:          "未登录",
		GatewayErrorCodeNotHavePermit:     "没有权限",
		GatewayErrorCodePermitConfigError: "网关配置错误(%s)",
		GateWayPathError:                  "访问路径异常(%s)",
		GateWayLoadConfigError:            "网关加载路由配置异常(%s)",
	}
)

type ResponseCallBack func()
