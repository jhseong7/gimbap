package salt

import (
	"github.com/jhseong7/gimbap/provider"
)

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

var SaltProvider = provider.DefineProvider(provider.ProviderOption{
	Name:         "SaltService",
	Instantiator: NewSalt,
})
