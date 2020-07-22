module kmodules.xyz/cert-manager-util

go 1.12

require (
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/jetstack/cert-manager v0.16.0-alpha.1
	github.com/pkg/errors v0.9.1
	k8s.io/apimachinery v0.18.5
	kmodules.xyz/client-go v0.0.0-20200521065424-173e32c78a20
)

replace (
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.19.0-alpha.0.0.20200520235721-10b58e57a423
	k8s.io/apiserver => github.com/kmodules/apiserver v0.18.4-0.20200521000930-14c5f6df9625
	k8s.io/kubernetes => github.com/kmodules/kubernetes v1.19.0-alpha.0.0.20200521033432-49d3646051ad
)
