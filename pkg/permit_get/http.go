package permit_get

import (
	"github.com/juetun/traefikplugins/pkg"
)

type HttpGet struct {
}

func (r *HttpGet) Do(uriParam *pkg.UriParam) (permitResult pkg.PermitGetResult, errCode int, errMsg string) {

	return
}

func NewHttpGet() (res pkg.PermitGet) {
	res = &HttpGet{}
	return
}
