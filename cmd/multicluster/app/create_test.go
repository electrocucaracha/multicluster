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

package app_test

import (
	"github.com/electrocucaracha/multicluster/cmd/multicluster/app"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func (m mock) Create(configPath, wanem string) error {
	return nil
}

var _ = Describe("Create Command", func() {
	var cmd *cobra.Command

	BeforeEach(func() {
		provider := mock{}
		cmd = app.NewCreateCommand(provider)
	})

	DescribeTable("creation execution process", func(shouldSucceed bool, args ...string) {
		cmd.SetArgs(args)
		err := cmd.Execute()

		if shouldSucceed {
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
		}
	},
		Entry("when the default options are provided", true),
		Entry("when the configuration path options are defined",
			true, "--config", "config.yml"),
		Entry("when the configuration path and WAN emulator image options are defined",
			true, "--config", "config.yml", "--wanem", "docker"),
		Entry("when an empty configuration path option is provided",
			false, "--name", "test", "--config", ""),
	)
})
