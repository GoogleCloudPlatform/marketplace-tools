package apply

type testResource struct {
	BaseResource
	applyFunc func(r Registry) error
	depFunc func() []Reference
}

func (tr *testResource) Apply(r Registry) error {
	return tr.applyFunc(r)
}

func (tr *testResource) GetDependencies() []Reference {
	if tr.depFunc == nil {
		return nil
	}
	return tr.depFunc()
}

func NewTestResource(name string) *testResource {
	return NewTestResourceFunc(name, nil, nil)
}

func NewTestResourceFunc(name string, applyFunc func(Registry) error, depFunc func() []Reference) *testResource {
	return &testResource{
		BaseResource: BaseResource{
			TypeMeta: TypeMeta{
				Kind:       "testKind",
				APIVersion: "testv1",
			},
			Metadata: Metadata{Name: name},
		},
		applyFunc: applyFunc,
		depFunc: depFunc,
	}
}