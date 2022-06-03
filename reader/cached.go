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

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	informers "github.com/cert-manager/cert-manager/pkg/client/informers/externalversions"
	listers "github.com/cert-manager/cert-manager/pkg/client/listers/certmanager/v1"
)

type cachedImpl struct {
	factory informers.SharedInformerFactory
	stopCh  <-chan struct{}

	lock         sync.RWMutex
	ciLister     listers.ClusterIssuerLister
	issuerLister listers.IssuerLister
}

var _ Reader = &cachedImpl{}

func (i *cachedImpl) ClusterIssuers() listers.ClusterIssuerLister {
	i.lock.RLock()
	if i.ciLister != nil {
		i.lock.RUnlock()
		return i.ciLister
	}
	i.lock.RUnlock()

	createLister := func() listers.ClusterIssuerLister {
		i.lock.Lock()
		defer i.lock.Unlock()
		if i.ciLister != nil {
			return i.ciLister
		}

		informerType := reflect.TypeOf(&cmapi.ClusterIssuer{})
		informerDep, _ := i.factory.ForResource(cmapi.SchemeGroupVersion.WithResource("clusterissuers"))
		i.factory.Start(i.stopCh)
		if synced := i.factory.WaitForCacheSync(i.stopCh); !synced[informerType] {
			panic(fmt.Sprintf("informer for %s hasn't synced", informerType))
		}
		i.ciLister = listers.NewClusterIssuerLister(informerDep.Informer().GetIndexer())
		return i.ciLister
	}
	return createLister()
}

func (i *cachedImpl) Issuers(namespace string) listers.IssuerNamespaceLister {
	i.lock.RLock()
	if i.issuerLister != nil {
		i.lock.RUnlock()
		return i.issuerLister.Issuers(namespace)
	}
	i.lock.RUnlock()

	createLister := func() listers.IssuerLister {
		i.lock.Lock()
		defer i.lock.Unlock()
		if i.issuerLister != nil {
			return i.issuerLister
		}

		informerType := reflect.TypeOf(&cmapi.Issuer{})
		informerDep, _ := i.factory.ForResource(cmapi.SchemeGroupVersion.WithResource("issuers"))
		i.factory.Start(i.stopCh)
		if synced := i.factory.WaitForCacheSync(i.stopCh); !synced[informerType] {
			panic(fmt.Sprintf("informer for %s hasn't synced", informerType))
		}
		i.issuerLister = listers.NewIssuerLister(informerDep.Informer().GetIndexer())
		return i.issuerLister
	}
	return createLister().Issuers(namespace)
}
