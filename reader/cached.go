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
	"fmt"
	"reflect"
	"sync"

	core "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	informers "github.com/jetstack/cert-manager/pkg/client/informers/externalversions"
	v1 "github.com/jetstack/cert-manager/pkg/client/listers/certmanager/v1"
)

type cachedImpl struct {
	factory informers.SharedInformerFactory
	stopCh  <-chan struct{}

	lock         sync.RWMutex
	ciLister     v1.ClusterIssuerLister
	issuerLister v1.IssuerLister
}

var _ ConfigReader = &cachedImpl{}

func (i *cachedImpl) ClusterIssuers() v1.ClusterIssuerLister {
	getLister := func() v1.ClusterIssuerLister {
		i.lock.RLock()
		defer i.lock.RUnlock()
		if i.ciLister != nil {
			return i.ciLister
		}

		informerType := reflect.TypeOf(&core.ClusterIssuer{})
		informerDep, _ := i.factory.ForResource(core.SchemeGroupVersion.WithResource("clusterissuers"))
		i.factory.Start(i.stopCh)
		if synced := i.factory.WaitForCacheSync(i.stopCh); !synced[informerType] {
			panic(fmt.Sprintf("informer for %s hasn't synced", informerType))
		}
		i.ciLister = v1.NewClusterIssuerLister(informerDep.Informer().GetIndexer())
		return i.ciLister
	}
	return getLister()
}

func (i *cachedImpl) Issuers(namespace string) v1.IssuerNamespaceLister {
	getLister := func() v1.IssuerLister {
		i.lock.RLock()
		defer i.lock.RUnlock()
		if i.issuerLister != nil {
			return i.issuerLister
		}

		informerType := reflect.TypeOf(&core.Issuer{})
		informerDep, _ := i.factory.ForResource(core.SchemeGroupVersion.WithResource("issuers"))
		i.factory.Start(i.stopCh)
		if synced := i.factory.WaitForCacheSync(i.stopCh); !synced[informerType] {
			panic(fmt.Sprintf("informer for %s hasn't synced", informerType))
		}
		i.issuerLister = v1.NewIssuerLister(informerDep.Informer().GetIndexer())
		return i.issuerLister
	}
	return getLister().Issuers(namespace)
}
