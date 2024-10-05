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

	// Multiple dependency return tester
	MultiA struct{}
	MultiB struct{}
	MultiC struct{}
	MultiD struct{}
	MultiE struct{}
	MultiF struct{}
	MultiG struct{}
	MultiH struct{}
)

// A --> B --> C
// D --> B
// E --> F --> B
func NewA() *A                 { return &A{} }
func NewB(a *A, d *D, f *F) *B { return &B{} }
func NewC(b *B) *C             { return &C{} }
func NewD() *D                 { return &D{} }
func NewE() *E                 { return &E{} }
func NewF(e *E) *F             { return &F{} }

// Circular dependency tester
// CirA --> CirB --> CirC --> CirA
func NewCirA(c *CirC) *CirA { return &CirA{} }
func NewCirB(a *CirA) *CirB { return &CirB{} }
func NewCirC(b *CirB) *CirC { return &CirC{} }

// Unresolved dependency tester
// OrphanA --> OrphanB
// But will be tested with only OrphanB + A,B,C,D,E,F
func NewOrphanA() *OrphanA           { return &OrphanA{} }
func NewOrphanB(a *OrphanA) *OprhanB { return &OprhanB{} }

// Multiple dependency return tester
// A single instantiator that returns multiple values
// ABC -> D, D -> EF, E --> G,  F --> H
func NewMultiABC() (*MultiA, *MultiB, *MultiC)          { return &MultiA{}, &MultiB{}, &MultiC{} }
func NewMultiD(a *MultiA, b *MultiB, c *MultiC) *MultiD { return &MultiD{} }
func NewMultiEF(d *MultiD) (*MultiE, *MultiF)           { return &MultiE{}, &MultiF{} }
func NewMultiG(e *MultiE) *MultiG                       { return &MultiG{} }
func NewMultiH(f *MultiF) *MultiH                       { return &MultiH{} }

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

var MultiABCProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "MultiABC", Instantiator: NewMultiABC})
var MultiDProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "MultiD", Instantiator: NewMultiD})
var MultiEFProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "MultiEF", Instantiator: NewMultiEF})
var MultiGProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "MultiG", Instantiator: NewMultiG})
var MultiHProvider = gimbap.DefineProvider(gimbap.ProviderOption{Name: "MultiH", Instantiator: NewMultiH})

var _ = Describe("GimbapDependencyManager", func() {

	Context("Test resolver", func() {
		fmt.Println("addProvider")

		It("Normal resolving", func() {
			instanceMap := make(map[reflect.Type]reflect.Value)

			providerList := []*provider.Provider{
				AProvider,
				BProvider,
				CProvider,
				DProvider,
				EProvider,
				FProvider,
			}
			gimbapManager := manager.NewGimbapDependencyManager()
			gimbapManager.ResolveDependencies(instanceMap, providerList)
		})

		It("Circular dependency --> will panic", func() {
			instanceMap := make(map[reflect.Type]reflect.Value)

			providerList := []*provider.Provider{
				CirAProvider,
				CirBProvider,
				CirCProvider,
			}

			Expect(func() {
				gimbapManager := manager.NewGimbapDependencyManager()
				gimbapManager.ResolveDependencies(instanceMap, providerList)
			}).To(Panic())
		})

		It("Orphan dependency --> will panic", func() {
			instanceMap := make(map[reflect.Type]reflect.Value)

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
				gimbapManager := manager.NewGimbapDependencyManager()
				gimbapManager.ResolveDependencies(instanceMap, providerList)
			}).To(Panic())
		})

		It("Multiple dependency return", func() {
			instanceMap := make(map[reflect.Type]reflect.Value)

			providerList := []*provider.Provider{
				MultiABCProvider,
				MultiDProvider,
				MultiEFProvider,
				MultiGProvider,
				MultiHProvider,
			}

			gimbapManager := manager.NewGimbapDependencyManager()
			gimbapManager.ResolveDependencies(instanceMap, providerList)

			Expect(instanceMap[reflect.TypeOf(&MultiA{})]).ToNot(BeNil())
		})

		It("Duplicate provider for 1 type provided (start node) --> will panic", func() {
			instanceMap := make(map[reflect.Type]reflect.Value)

			providerList := []*provider.Provider{
				AProvider,
				AProvider, // A is provided twice
				BProvider,
				CProvider,
				DProvider,
				EProvider,
				FProvider,
			}

			Expect(func() {
				gimbapManager := manager.NewGimbapDependencyManager()
				gimbapManager.ResolveDependencies(instanceMap, providerList)
			}).To(Panic())
		})

		It("Duplicate provider for 1 type provided (node) --> will panic", func() {
			instanceMap := make(map[reflect.Type]reflect.Value)

			providerList := []*provider.Provider{
				AProvider,
				BProvider,
				BProvider, // B is provided twice
				CProvider,
				DProvider,
				EProvider,
				FProvider,
			}

			Expect(func() {
				gimbapManager := manager.NewGimbapDependencyManager()
				gimbapManager.ResolveDependencies(instanceMap, providerList)
			}).To(Panic())
		})
	})
})
