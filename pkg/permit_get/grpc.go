package permit_get

import (
	"github.com/juetun/traefikplugins/pkg"
)

type GrpcGet struct {
}

func (r *GrpcGet) Do(uriParam *pkg.UriParam) (permitResult pkg.PermitGetResult, errCode int, errMsg string) {

	return
}

func NewGrpcGet() pkg.PermitGet {
	p := &GrpcGet{}
	return p
}
