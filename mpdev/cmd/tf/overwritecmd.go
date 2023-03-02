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

package tf

import (
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/docs"
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/tf"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// GetOverwriteCommand returns `overwrite` command used to create mpdev resources.
func GetOverwriteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "overwrite",
		Short:   docs.OverwriteShort,
		Long:    docs.OverwriteLong,
		Example: docs.OverwriteExamples,
		RunE:    overwriteRunE,
	}
	cmd.SilenceUsage = true

	return cmd
}

func overwriteRunE(_ *cobra.Command, _ []string) (err error) {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	err = tf.Overwrite(bytes, dir)
	return err
}
