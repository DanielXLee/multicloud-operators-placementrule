// Copyright 2021 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitopscluster

import (
	stdlog "log"
	"os"
	"path/filepath"
	"sync"
	"testing"

	spokeClusterV1 "github.com/clusternet/clusternet/pkg/apis/clusters/v1beta1"
	"github.com/onsi/gomega"
	clusterv1alpha1 "github.com/open-cluster-management/api/cluster/v1alpha1"
	endpointapis "github.com/open-cluster-management/klusterlet-addon-controller/pkg/apis"
	gitopsclusterV1alpha1 "github.com/open-cluster-management/multicloud-operators-placementrule/pkg/apis/apps/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var cfg *rest.Config

func TestMain(m *testing.M) {
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "deploy", "crds"),
			filepath.Join("..", "..", "..", "hack", "test", "crds"),
		},
	}

	endpointapis.AddToScheme(scheme.Scheme)
	spokeClusterV1.AddToScheme(scheme.Scheme)
	clusterv1alpha1.AddToScheme(scheme.Scheme)

	gitopsclusterV1alpha1.AddToScheme(scheme.Scheme)

	var err error
	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}

	code := m.Run()

	t.Stop()

	os.Exit(code)
}

// SetupTestReconcile returns a reconcile.Reconcile implementation that delegates to inner
func SetupTestReconcile(inner reconcile.Reconciler) reconcile.Reconciler {
	fn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(req)

		return result, err
	})

	return fn
}

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager, g *gomega.GomegaWithT) (chan struct{}, *sync.WaitGroup) {
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(stop)).NotTo(gomega.HaveOccurred())
	}()

	return stop, wg
}
