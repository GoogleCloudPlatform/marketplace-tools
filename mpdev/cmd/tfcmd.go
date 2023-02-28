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
