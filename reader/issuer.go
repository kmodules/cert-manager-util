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

	core "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	cmcs "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	v1 "github.com/jetstack/cert-manager/pkg/client/listers/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/pager"
)

var _ v1.IssuerNamespaceLister = &issuerNamespaceLister{}

// issuerNamespaceLister implements the NamespaceLister interface.
type issuerNamespaceLister struct {
	dc        cmcs.Interface
	namespace string
}

// List lists all resources in the indexer for a given namespace.
func (l *issuerNamespaceLister) List(selector labels.Selector) (ret []*core.Issuer, err error) {
	fn := func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		return l.dc.CertmanagerV1().Issuers(l.namespace).List(ctx, opts)
	}
	opts := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	err = pager.New(fn).EachListItem(context.TODO(), opts, func(obj runtime.Object) error {
		o, ok := obj.(*core.Issuer)
		if !ok {
			return fmt.Errorf("expected *core.Issuer, found %s", reflect.TypeOf(obj))
		}
		ret = append(ret, o)
		return nil
	})
	return ret, err
}

// Get retrieves a resource from the indexer for a given namespace and name.
func (l *issuerNamespaceLister) Get(name string) (*core.Issuer, error) {
	obj, err := l.dc.CertmanagerV1().Issuers(l.namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}
