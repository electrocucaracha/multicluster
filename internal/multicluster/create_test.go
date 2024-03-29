/*
Copyright © 2023

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

package multicluster_test

import (
	"github.com/electrocucaracha/multicluster/internal/multicluster"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"
)

type mockContainerProvider struct {
	mockBase
}

func (m *mockContainerProvider) CreateNetwork(string, string, bool) error {
	return m.popError()
}

func (m *mockContainerProvider) ConnectNetwork(string, string, string) error {
	return m.popError()
}

func (m *mockContainerProvider) ReplaceGateway(string, string) error {
	return m.popError()
}

func (m *mockContainerProvider) ListNetwork() ([]string, error) {
	return nil, m.popError()
}

func (m *mockWanProvider) Create(string, string) (string, error) {
	return "", m.popError()
}

func (m *mockWanProvider) AddRoutes(string, string, ...string) error {
	return m.popError()
}

func (m *mockClusterProvider) Create(string, ...cluster.CreateOption) error {
	return m.popError()
}

func (m *mockClusterProvider) ListNodes(string) ([]multicluster.Node, error) {
	return m.Nodes, m.popError()
}

var _ = Describe("Create Service", func() {
	var provider *multicluster.KindDataSource
	var clusterProvider *mockClusterProvider
	var wanProvider *mockWanProvider
	var configReader *mockConfigReader
	var containerProvider *mockContainerProvider
	emptyConfig := multicluster.Config{}
	testConfig := multicluster.Config{
		Name: "Lab",
		Clusters: map[string]multicluster.ClusterConfig{
			"test": {
				NodeSubnet: "172.88.0.0/16",
				Cluster: &v1alpha4.Cluster{
					Networking: v1alpha4.Networking{
						PodSubnet:     "10.196.0.0/16",
						ServiceSubnet: "10.96.0.0/16",
					},
				},
			},
		},
	}

	BeforeEach(func() {
		logger := log.NoopLogger{}
		clusterProvider = &mockClusterProvider{}
		wanProvider = &mockWanProvider{}
		configReader = &mockConfigReader{}
		containerProvider = &mockContainerProvider{}

		provider = multicluster.NewProvider(configReader, wanProvider,
			clusterProvider, containerProvider, logger)
	})

	DescribeTable("create execution service process", func(
		config multicluster.Config, clusters []string,
		wanErrorMessages []string, clusterProviderErrorMessages []string,
		containerProviderErrorMessages []string, shouldSucceed bool,
	) {
		configReader.ConfigInfo = config
		clusterProvider.Clusters = clusters
		errMsgExpected := wanProvider.PushErrorMessages(wanErrorMessages)
		if errMsgExpected == "" {
			errMsgExpected = clusterProvider.PushErrorMessages(clusterProviderErrorMessages)
		}
		if errMsgExpected == "" {
			errMsgExpected = containerProvider.PushErrorMessages(containerProviderErrorMessages)
		}

		err := provider.Create("configPath", "wanem")
		if shouldSucceed {
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(errMsgExpected))
		}
	},
		Entry("when empty cluster config is provided",
			emptyConfig, []string{""}, nil, nil, nil, true),
		Entry("when a valid cluster config is provided",
			testConfig, []string{"node01"}, nil, nil, nil, true),
	)
})
