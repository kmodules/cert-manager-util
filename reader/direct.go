/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reader

import (
	"sync"

	cmcs "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	v1 "github.com/jetstack/cert-manager/pkg/client/listers/certmanager/v1"
)

type directImpl struct {
	dc cmcs.Interface

	lock         sync.RWMutex
	ciLister     v1.ClusterIssuerLister
	issuerLister v1.IssuerNamespaceLister
}

var _ Reader = &directImpl{}

func (i *directImpl) ClusterIssuers() v1.ClusterIssuerLister {
	i.lock.RLock()
	defer i.lock.RUnlock()
	if i.ciLister != nil {
		return i.ciLister
	}

	i.ciLister = &clusterissuerLister{dc: i.dc}
	return i.ciLister
}

func (i *directImpl) Issuers(namespace string) v1.IssuerNamespaceLister {
	i.lock.RLock()
	defer i.lock.RUnlock()
	if i.issuerLister != nil {
		return i.issuerLister
	}

	i.issuerLister = &issuerNamespaceLister{dc: i.dc, namespace: namespace}
	return i.issuerLister
}
