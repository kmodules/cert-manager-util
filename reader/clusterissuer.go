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
	"context"
	"fmt"
	"reflect"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmcs "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	listers "github.com/cert-manager/cert-manager/pkg/client/listers/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/pager"
)

var _ listers.ClusterIssuerLister = &clusterissuerLister{}

// clusterissuerLister implements the NamespaceLister interface.
type clusterissuerLister struct {
	dc cmcs.Interface
}

// List lists all resources in the indexer.
func (l *clusterissuerLister) List(selector labels.Selector) (ret []*cmapi.ClusterIssuer, err error) {
	fn := func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		return l.dc.CertmanagerV1().ClusterIssuers().List(ctx, opts)
	}
	opts := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	err = pager.New(fn).EachListItem(context.TODO(), opts, func(obj runtime.Object) error {
		o, ok := obj.(*cmapi.ClusterIssuer)
		if !ok {
			return fmt.Errorf("expected *cmapi.ClusterIssuer, found %s", reflect.TypeOf(obj))
		}
		ret = append(ret, o)
		return nil
	})
	return ret, err
}

// Get retrieves a resource from the indexer.
func (l *clusterissuerLister) Get(name string) (*cmapi.ClusterIssuer, error) {
	obj, err := l.dc.CertmanagerV1().ClusterIssuers().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}
