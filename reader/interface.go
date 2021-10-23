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
	"time"

	cmcs "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	informers "github.com/jetstack/cert-manager/pkg/client/informers/externalversions"
	v1 "github.com/jetstack/cert-manager/pkg/client/listers/certmanager/v1"
)

type ConfigReader interface {
	ClusterIssuers() v1.ClusterIssuerLister
	Issuers(namespace string) v1.IssuerNamespaceLister
}

func New(dc cmcs.Interface) ConfigReader {
	return &directImpl{
		dc: dc,
	}
}

func NewCached(dc cmcs.Interface, defaultResync time.Duration, stopCh <-chan struct{}) ConfigReader {
	return &cachedImpl{
		factory: informers.NewSharedInformerFactory(dc, defaultResync),
		stopCh:  stopCh,
	}
}

func NewCachedWithOptions(dc cmcs.Interface, defaultResync time.Duration, stopCh <-chan struct{}, options ...informers.SharedInformerOption) ConfigReader {
	return &cachedImpl{
		factory: informers.NewSharedInformerFactoryWithOptions(dc, defaultResync, options...),
		stopCh:  stopCh,
	}
}

func NewSharedCached(factory informers.SharedInformerFactory, stopCh <-chan struct{}) ConfigReader {
	return &cachedImpl{
		factory: factory,
		stopCh:  stopCh,
	}
}
