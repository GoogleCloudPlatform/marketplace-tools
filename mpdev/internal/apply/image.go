package apply

type PackerGceImageBuilder struct {
	TypeMeta
	ObjectMeta
	Builder struct {
		Script struct {
			File string
		}
	}

	Tests []struct {
		Name   string
		Script struct {
			File string
		}
	}
}

func (p *PackerGceImageBuilder) Apply() {

}

type GceImage struct {
	TypeMeta
	ObjectMeta `json:"metadata"`
	ImageRef   Reference
	BuilderRef Reference
	Image      Image
}

func (g *GceImage) Apply() {

}

type Image struct {
	ProjectId          string
	NamePartsSeparator string
	NameParts          []string
}