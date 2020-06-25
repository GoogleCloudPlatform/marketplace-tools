package apply

type PackerGceImageBuilder struct {
	ResourceShared
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

func (p *PackerGceImageBuilder) Apply() error {
	return nil
}

type GceImage struct {
	ResourceShared
	ImageRef   Reference
	BuilderRef Reference
	Image      Image
}

func (g *GceImage) Apply() error {
	return nil
}

type Image struct {
	ProjectId          string
	NamePartsSeparator string
	NameParts          []string
}