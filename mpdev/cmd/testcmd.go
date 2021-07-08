package cmd

import (
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/docs"
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/resources"
	"github.com/spf13/cobra"
	"k8s.io/utils/exec"
)

// GetTestCommand returns `test` command used to create mpdev resources.
func GetTestCommand() *cobra.Command {
	var c command
	cmd := &cobra.Command{
		Use:     "test -f FILENAME [--dryRun]",
		Short:   docs.TestShort,
		Long:    docs.TestLong,
		Example: docs.TestExamples,
		RunE:    c.testRunE,
	}

	cmd.Flags().BoolVar(&c.DryRun, "dryrun", c.DryRun, "if set, validates configuration files without creating resource")
	cmd.Flags().StringSliceVarP(&c.Filenames, "filename", "f", c.Filenames, "that contains the configuration to test")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "filename")

	return cmd
}

// RunE Executes the `test` command
func (c *command) testRunE(_ *cobra.Command, _ []string) (err error) {
	registry := resources.NewRegistry(exec.New())
	err = resources.PopulateRegistryFromFiles(registry, c.Filenames)
	if err != nil {
		return err
	}

	err = registry.Test(c.DryRun)
	return err
}
