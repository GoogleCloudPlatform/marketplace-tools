// Copyright 2020 Google LLC
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

package cmd

import (
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/cmd/tf"
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/docs"
	"github.com/spf13/cobra"
)

// GetTfCommand returns `tf` command used to create mpdev resources.
func GetTfCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tf",
		Short: docs.TfShort,
		Long:  docs.TfLong,
	}

	cmd.AddCommand(tf.GetOverwriteCommand())
	return cmd
}
