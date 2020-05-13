// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"path/filepath"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/docs"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/cmd/config/ext"
)

const kptFileName = "Kptfile"

func GetMain() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mpdev",
		Short:   docs.ReferenceShort,
		Long:    docs.ReferenceLong,
		Example: docs.ReferenceExamples,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Override openApi file location such that KptFile will be modified
			// by mpdev cfg commands.
			// See: https://github.com/GoogleContainerTools/kpt/blob/bf211c225274fe6747304c9b6bf55ea5a98b603a/run/run.go#L48
			ext.GetOpenAPIFile = func(args []string) (s string, err error) {
				return filepath.Join(args[0], kptFileName), nil
			}
		},
	}
	cmd.AddCommand(GetMpdevCommands("mpdev")...)

	return cmd
}
