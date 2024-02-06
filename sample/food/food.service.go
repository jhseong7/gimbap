package food

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/jhseong7/nassi-golang/core"
	"github.com/jhseong7/nassi-golang/sample/salt"
)

type (
	FoodService struct {
		Name   string
		Salt   salt.SaltService
		random float64
	}
)

func (f *FoodService) GetName() string {
	return f.Name + " with " + f.Salt.GetName() + " and random " + fmt.Sprintf("%f", f.random)
}

func NewFood(salt *salt.SaltService) *FoodService {
	return &FoodService{
		Name:   "Food",
		Salt:   *salt,
		random: math.Floor(1000*rand.Float64()) + 1,
	}
}

var FoodProvider = core.DefineProvider(core.ProviderOption{
	Name:         "FoodService",
	Instantiator: NewFood,
})
