package apply

import "strings"

const API_VERSION = "dev.marketplace.cloud.google.com/v1alpha1"
var TypeMapper = map[TypeMeta]func() Resource{
	{ ApiVersion: API_VERSION, Kind: "GceImage" } :             func() Resource { return &GceImage{} },
	{ ApiVersion: API_VERSION, Kind: "PackerGceImageBuilder"} : func() Resource { return &PackerGceImageBuilder{} },
	{ ApiVersion: API_VERSION, Kind: "DeploymentManagerAutogenTemplate"} : func() Resource { return &DeploymentManagerAutogenTemplate{}},
	{ ApiVersion: API_VERSION, Kind: "DeploymentManagerTemplateOnGCS"} : func() Resource { return &DeploymentManagerTemplateOnGCS{}},
}

type TypeMeta struct {
	Kind       string
	ApiVersion string
}

type Metadata struct {
	Name        string
	Annotations map[string]string
}

func (m *Metadata) GetName() string {
	return m.Name
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
	GetReference() Reference
	SetReferenceMap(ReferenceMap)
}

type ResourceShared struct {
	TypeMeta
	Metadata Metadata

	referenceMap ReferenceMap
}

func (rs *ResourceShared) GetReference() Reference {
	groupAndVersion := strings.Split(rs.ApiVersion, "/")
	return Reference{
		Group:   groupAndVersion[0],
		Kind:    rs.Kind,
		Name:    rs.Metadata.Name,
	}
}

func (rs *ResourceShared) SetReferenceMap(referenceMap ReferenceMap) {
	rs.referenceMap = referenceMap
}

type Reference struct {
	Group string
	Kind  string
	Name  string
}

type ReferenceMap map[Reference]Resource
