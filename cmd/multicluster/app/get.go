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

package app

import (
	"github.com/electrocucaracha/multicluster/internal/multicluster"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
)

func NewGetCommand(provider multicluster.DataSource) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the clusters that belong to the multi cluster",
		Long:  `Get the clusters that belong to the multi cluster`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := getName(cmd.Flags())
			if err != nil {
				return errors.Wrap(err, "failed to retrieve the name of the multi-cluster")
			}

			if err := provider.Get(name); err != nil {
				return errors.Wrapf(err, "failed to retrieve %s multi-cluster info", name)
			}

			return nil
		},
	}

	cmd.Flags().String(
		"name",
		cluster.DefaultName,
		"the multi-cluster context name",
	)

	return cmd
}
