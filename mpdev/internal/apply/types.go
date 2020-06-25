package apply


const API_VERSION = "dev.marketplace.cloud.google.com/v1alpha1"
var TypeMapper = map[TypeMeta]func() Resource{
	{ ApiVersion: API_VERSION, Kind: "GceImage" } :             func() Resource { return &GceImage{} },
	{ ApiVersion: API_VERSION, Kind: "PackerGceImageBuilder"} : func() Resource { return &PackerGceImageBuilder{} },
	{ ApiVersion: API_VERSION, Kind: "DeploymentManagerAutogenTemplate"} : func() Resource { return &DeploymentManagerAutogenTemplate{}},
}

type TypeMeta struct {
	Kind       string
	ApiVersion string
}

type Metadata struct {
	Name        string
	Annotations map[string]string
}

type Unstructured map[string]interface{}

func (u *Unstructured) GetTypeMeta() TypeMeta {
	kind, ok := (*u)["kind"].(string)
	if !ok {
		return TypeMeta{}
	}

	apiVersion, ok := (*u)["apiVersion"].(string)
	if !ok {
		return TypeMeta{}
	}
	return TypeMeta{
		Kind:       kind,
		ApiVersion: apiVersion,
	}
}


type Resource interface {
	Apply() error
}

type Reference struct {
	Group string
	Kind  string
	Name  string
}
