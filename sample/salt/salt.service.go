package salt

import "github.com/jhseong7/nassi-golang/core"

type (
	SaltService struct {
		Name string
	}
)

func (s *SaltService) GetName() string {
	return s.Name
}

func NewSalt() *SaltService {
	return &SaltService{
		Name: "Salty Salt!!",
	}
}

var SaltProvider = core.DefineProvider(core.ProviderOption{
	Name:         "SaltService",
	Instantiator: NewSalt,
})
