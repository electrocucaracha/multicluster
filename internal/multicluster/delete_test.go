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
	"sigs.k8s.io/kind/pkg/log"
)

func (m *mockClusterProvider) Delete(string, string) error {
	return m.popError()
}

func (m *mockWanProvider) Delete(string) error {
	return m.popError()
}

func (m *mockContainerProvider) DeleteNetwork(name string) error {
	return m.popError()
}

var _ = Describe("Delete Service", func() {
	var provider *multicluster.KindDataSource
	var clusterProvider *mockClusterProvider
	var wanProvider *mockWanProvider
	var configReader *mockConfigReader
	var containerProvider *mockContainerProvider
	emptyConfig := multicluster.Config{}
	testConfig := multicluster.Config{
		Name: "testlab",
		Clusters: map[string]multicluster.ClusterConfig{
			"test": {},
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

	DescribeTable("delete execution service process", func(
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

		err := provider.Delete("name", "configPath")
		if shouldSucceed {
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(errMsgExpected))
		}
	},
		Entry("when empty cluster config is provided",
			emptyConfig, []string{""}, nil, nil, nil, true),
		Entry("when non-empty cluster config is provided",
			testConfig, []string{""}, nil, nil, nil, true),
		Entry("when non-empty cluster config is provided but no cluster matches",
			testConfig, []string{"kind"}, nil, nil, nil, true),
		Entry("when non-empty cluster config is provided and cluster matches",
			testConfig, []string{"test"}, nil, nil, nil, true),
		Entry("when WAN emulator raises an error during deletion",
			emptyConfig, []string{""}, []string{"wan provider error"}, nil,
			nil, false),
		Entry("when cluster provider raises an error during retrieval",
			emptyConfig, []string{""}, nil,
			[]string{"cluster provider retrieval error"}, nil, false),
		Entry("when cluster provider raises an error during deletion",
			testConfig, []string{"test"}, nil,
			[]string{"cluster provider deletion error", ""}, nil, true),
		Entry("when container provider raises an error during retrieval",
			testConfig, []string{"test"}, nil, nil,
			[]string{"container provider retrieval error"}, false),
		Entry("when container provider raises an error during deletion",
			testConfig, []string{"test"}, nil, nil,
			[]string{"container provider deletion error", ""}, true),
	)
})
