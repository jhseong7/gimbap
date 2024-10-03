package dependency_test

import (
	"fmt"
	"reflect"

	"github.com/jhseong7/gimbap"
	manager "github.com/jhseong7/gimbap/dependency"
	"github.com/jhseong7/gimbap/provider"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/**
A --> B --> C
D --> B
E --> F --> B
*/

// Structs to test the dependency manager
type (
	A struct{}
	B struct{}
	C struct{}
	D struct{}
	E struct{}
	F struct{}
	G struct{}

	// Circular dependency tester
	CirA struct{}
	CirB struct{}
	CirC struct{}

	// Unresolved dependency tester
	OrphanA struct{}
	OprhanB struct{}
)

func NewA() *A                 { return &A{} }
func NewB(a *A, d *D, f *F) *B { return &B{} }
func NewC(b *B) *C             { return &C{} }
func NewD() *D                 { return &D{} }
func NewE() *E                 { return &E{} }
func NewF(e *E) *F             { return &F{} }

func NewCirA(c *CirC) *CirA { return &CirA{} }
func NewCirB(a *CirA) *CirB { return &CirB{} }
func NewCirC(b *CirB) *CirC { return &CirC{} }

func NewOrphanA() *OrphanA           { return &OrphanA{} }
func NewOrphanB(a *OrphanA) *OprhanB { return &OprhanB{} }

var AProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "A", Instantiator: NewA})
var BProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "B", Instantiator: NewB})
var CProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "C", Instantiator: NewC})
var DProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "D", Instantiator: NewD})
var EProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "E", Instantiator: NewE})
var FProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "F", Instantiator: NewF})

var CirAProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "CirA", Instantiator: NewCirA})
var CirBProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "CirB", Instantiator: NewCirB})
var CirCProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "CirC", Instantiator: NewCirC})

var OrphanAProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "OrphanA", Instantiator: NewOrphanA})
var OrphanBProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "OrphanB", Instantiator: NewOrphanB})

var _ = Describe("GimbapDependencyManager", func() {
	gimbapManager := manager.NewGimbapDependencyManager()

	Context("Test resolver", func() {
		fmt.Println("addProvider")

		instanceMap := make(map[reflect.Type]reflect.Value)

		It("Normal resolving", func() {
			providerList := []*provider.Provider{
				AProvider,
				BProvider,
				CProvider,
				DProvider,
				EProvider,
				FProvider,
			}

			gimbapManager.ResolveDependencies(instanceMap, providerList)
		})

		It("Circular dependency --> will panic", func() {
			providerList := []*provider.Provider{
				CirAProvider,
				CirBProvider,
				CirCProvider,
			}

			Expect(func() {
				gimbapManager.ResolveDependencies(instanceMap, providerList)
			}).To(Panic())
		})

		It("Orphan dependency --> will panic", func() {
			providerList := []*provider.Provider{
				AProvider,
				BProvider,
				CProvider,
				DProvider,
				EProvider,
				FProvider,
				OrphanBProvider, // Only add B to trigger the not resolved error
			}

			Expect(func() {
				gimbapManager.ResolveDependencies(instanceMap, providerList)
			}).To(Panic())
		})
	})
})
