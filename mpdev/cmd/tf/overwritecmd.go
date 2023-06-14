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

	config, err := tf.GetOverwriteConfig(bytes)
	if err != nil {
		return err
	}

	err = tf.OverwriteTf(config, dir)
	if err != nil {
		return err
	}

	return tf.OverwriteMetadata(config, dir)
}
