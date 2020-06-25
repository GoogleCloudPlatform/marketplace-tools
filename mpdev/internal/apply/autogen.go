package apply

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

type DeploymentManagerAutogenTemplate struct {
	ResourceShared
	AutogenFile string
	PartnerId string
	SolutionId string

	OutDir string
}

type ContainerProcess struct {
	containerImage string
	processArgs []string
	mounts []string
}

func (dm *DeploymentManagerAutogenTemplate) Apply() error {
	// Name from running `docker build mpdev/autogen -f mpdev/autogen/Dockerfile -t autogen_converter
	autogenConverterImg := "autogen_converter"
	cmd := getCommand(ContainerProcess{
		containerImage: autogenConverterImg,
		processArgs:    []string{"--partnerId", dm.PartnerId, "--solutionId", dm.SolutionId},
		mounts:         nil,
	})

	f, err := os.Open(dm.AutogenFile)
	if err != nil {
		return err
	}
	defer f.Close()

	dir, err := ioutil.TempDir("", "autogen")
	if err != nil {
		return err
	}

	autogenFile, err := os.Create(filepath.Join(dir, "autogen.yaml"))
	if err != nil {
		return err
	}
	defer autogenFile.Close()

	cmd.Stderr = os.Stderr
	cmd.Stdin = f
	cmd.Stdout = autogenFile

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute autogen converter")
	}

	autogenImg := "gcr.io/cloud-marketplace-tools/dm/autogen"
	args := []string{"--input_type", "YAML", "--single_input", "/autogen/autogen.yaml",
		"--output_type", "PACKAGE", "--output", "/autogen"}

	cmd = getCommand(ContainerProcess{
		containerImage: autogenImg,
		processArgs:    args,
		mounts:         []string{fmt.Sprintf("type=bind,src=%s,dst=/autogen", dir)},
	})
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute autogen")
	}

	fmt.Printf("Wrote autogen output to directory: %s\n", dir)
	dm.OutDir = dir

	return nil
}

func getCommand(process ContainerProcess) *exec.Cmd {
	args := []string{"docker", "run", "--rm", "-i" }
	for _, mount := range process.mounts {
		args = append(args, "--mount", mount)
	}
	args = append(args, process.containerImage)
	args = append(args, process.processArgs...)
	return exec.Command(args[0], args[1:]...)
}

type DeploymentManagerTemplateOnGCS struct {
	ResourceShared
	AutogenRef Reference
	GCS struct {
		Bucket string
		Object string
	}
}

func (dm *DeploymentManagerTemplateOnGCS) Apply() error {
	autogenRef := dm.referenceMap[dm.AutogenRef]
	if autogenRef == nil {
		return fmt.Errorf("autogen template not found %+v", dm.AutogenRef)
	}

	autogenTemplate, ok := autogenRef.(*DeploymentManagerAutogenTemplate)
	if !ok {
		return fmt.Errorf("referenced autogen template is not correct type %+v", dm.AutogenRef)
	}

	gcsPath := fmt.Sprintf("gs://%s/%s", dm.GCS.Bucket, dm.GCS.Object)

	cmd := exec.Command("gsutil", "cp", "-r", autogenTemplate.OutDir, gcsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Uploading DM template to GCS. Running command: %v\n", cmd)

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to copy autogen template to GCS")
	}

	fmt.Printf("Uploaded DM template to GCS path: %s\n", gcsPath)

	return nil
}