module github.com/GoogleCloudPlatform/marketplace-tools

go 1.23

require (
	github.com/GoogleContainerTools/kpt v0.33.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/hcl/v2 v2.16.1
	github.com/hashicorp/terraform-config-inspect v0.0.0-20230223165911-2d94e3d51111
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.1
	github.com/tidwall/gjson v1.14.2
	github.com/tidwall/sjson v1.2.5
	gonum.org/v1/gonum v0.7.0
	gopkg.in/yaml.v3 v3.0.0
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
	sigs.k8s.io/kustomize/cmd/config v0.6.0
	sigs.k8s.io/yaml v1.2.0
	golang.org/x/crypto v0.31.0 // indirect
)
