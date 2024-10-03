// Experimental Gimbap Dependency Manager
package dependency

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/provider"
	"github.com/jhseong7/gimbap/util"
)

type (
	GimbapDependencyManager struct {
		IDependencyManager
		logger ecl.Logger
	}

	// Context holder to make the dependency manager stateless
	GimbapDependencyManagerContext struct {
		/*
			Dependency graph to resolve the dependencies
			The map will register a list of nodes that require the key node
			The node will be consumed upon initialization

			If any node is not consumed then it must throw --> circular dependency or missing dependency
		*/
		DependencyGraph map[reflect.Type]map[reflect.Type]*dependencyNode

		/*
			The starting node list to resolve the dependencies
			These are nodes that do not require any other nodes to initialize
		*/
		StartNodeList []*dependencyNode

		// Reference to the instance map
		InstanceMap map[reflect.Type]reflect.Value
	}

	// Internal types to handle the dependency graph
	dependencyNode struct {
		nodeType     reflect.Type // Same as the key of the map
		requires     []reflect.Type
		provider     *provider.Provider
		requeueCount int // This requeue count cannot exceed the number of dependencies --> used to detect circular dependencies
	}
)

func NewGimbapDependencyManagerContext(instanceMap map[reflect.Type]reflect.Value) *GimbapDependencyManagerContext {
	return &GimbapDependencyManagerContext{
		DependencyGraph: make(map[reflect.Type]map[reflect.Type]*dependencyNode),
		InstanceMap:     instanceMap,
	}
}

func (g *GimbapDependencyManager) createInstanceFromInstantiator(instantiator interface{}, instanceMap map[reflect.Type]reflect.Value) (reflect.Value, bool, error) {
	// Get the input types of the instantiator
	inputTypes, ok := util.DeriveInputTypesFromInstantiator(instantiator)
	if !ok {
		g.logger.Error("Failed to derive input types from instantiator. is the instantiator a function?")
		return reflect.Value{}, false, fmt.Errorf("failed to derive input types from instantiator")
	}

	// Get the input values from the instance map
	inputValues := make([]reflect.Value, len(inputTypes))
	for i, inputType := range inputTypes {
		inputValue, ok := instanceMap[inputType]
		if !ok { // The instance is not yet created --> return later
			return reflect.Value{}, false, nil
		}

		inputValues[i] = inputValue
	}

	// Call the instantiator with the input values
	returnValues := reflect.ValueOf(instantiator).Call(inputValues)

	// If the return values are empty --> panic
	if len(returnValues) == 0 {
		return reflect.Value{}, false, fmt.Errorf("failed to create instance from instantiator")
	}

	g.logger.Debugf("Provider instance created: %s", returnValues[0].Type().String())

	// Return the first return value
	return returnValues[0], true, nil
}

// Throw a panic + log any unresolved dependencies
// Build a readable error message for the user
func (g *GimbapDependencyManager) throwDependencyResolveError(context *GimbapDependencyManagerContext, msg string) {
	errorMessage := fmt.Sprintf("Failed to resolve the dependencies for the reason: %s\n\n", msg)

	unresolvedTypes := make(map[reflect.Type]*dependencyNode)

	// For all the nodes that are not resolved
	for _, nodeMap := range context.DependencyGraph {
		for _, node := range nodeMap {
			unresolvedTypes[node.nodeType] = node
		}
	}

	// For the unresolved types create a message
	errorMessage += "The dependency encountered unresolved dependencies:\n\n"
	for _, node := range unresolvedTypes {
		errorMessage += "( "
		for _, required := range node.requires {
			// Check if the required type is already resolved
			if _, ok := context.InstanceMap[required]; ok {
				continue
			}

			errorMessage += fmt.Sprintf("%v, ", required)
		}

		// Remove the last ", "
		errorMessage = errorMessage[:len(errorMessage)-2]

		errorMessage += fmt.Sprintf(" ) --> %v\n", node.provider.Name)
	}

	errorMessage += "\n\nPlease check the provider configuration"

	g.logger.Panic(errorMessage)
}

// Instantiate the providers using the node data
func (g *GimbapDependencyManager) instantiateProviders(context *GimbapDependencyManagerContext) {
	searchQueue := append([]*dependencyNode{}, context.StartNodeList...)
	inQueue := make(map[reflect.Type]bool)

	// While the search queue is not empty
	for len(searchQueue) > 0 {
		// Pop the first element
		node := searchQueue[0]
		searchQueue = searchQueue[1:]
		inQueue[node.nodeType] = false

		// Instantiate the provider
		inst, success, err := g.createInstanceFromInstantiator(node.provider.Instantiator, context.InstanceMap)
		if err != nil {
			g.logger.Panicf("Failed to instantiate provider: %s. Please check the provider configuration", node.provider.Name)
		}

		// The instance is not ready to be created --> add it to the search queue (will be resolved later)
		if !success {
			// Increase the requeue count
			node.requeueCount++
			// If the requeue count exceeds the number of dependencies --> circular dependency
			if node.requeueCount > len(node.requires) {
				g.throwDependencyResolveError(context, "Circular dependency detected")
			}

			searchQueue = append(searchQueue, node)
			inQueue[node.nodeType] = true
			continue
		}

		// Add the instance to the instance map
		context.InstanceMap[node.nodeType] = inst

		// Remove from the dependency graph if the instance is created
		for _, required := range node.requires {
			delete(context.DependencyGraph[required], node.nodeType)

			// Remove the root key if there are no more dependencies
			if len(context.DependencyGraph[required]) == 0 {
				delete(context.DependencyGraph, required)
			}
		}

		// Get the nodes that require the current node
		nodeMap, ok := context.DependencyGraph[node.nodeType]
		if !ok { // If there are no nodes that require the current node
			continue
		}

		// For all nodes that require the current node --> push to the search queue
		for _, n := range nodeMap {
			// If the node is already in the queue --> skip
			if inQueue[n.nodeType] {
				continue
			}

			searchQueue = append(searchQueue, n)
			inQueue[n.nodeType] = true
		}
	}

	// Finally check if there are any nodes that are not resolved
	if len(context.DependencyGraph) > 0 {
		g.throwDependencyResolveError(context, "Failed to resolve the dependencies. Missing dependency detected")
	}
}

func (g *GimbapDependencyManager) addProvider(context *GimbapDependencyManagerContext, p *provider.Provider) {
	// Get the return type and input types of the instantiator
	returnType, ok := util.DeriveTypeFromInstantiator(p.Instantiator)

	if !ok {
		g.logger.Panicf("Failed to derive type from instantiator: %s", p.Name)
	}

	// Get the input types of the instantiator
	inputTypes, ok := util.DeriveInputTypesFromInstantiator(p.Instantiator)

	if !ok {
		g.logger.Panicf("Failed to derive input types from instantiator: %s", p.Name)
	}

	node := &dependencyNode{
		nodeType: returnType,
		requires: inputTypes,
		provider: p,
	}

	// For all the input types, add the node to the dependency graph
	if len(inputTypes) > 0 {
		for _, inputType := range inputTypes {
			// Initialize the dependency graph if it does not exist
			if _, ok := context.DependencyGraph[inputType]; !ok {
				context.DependencyGraph[inputType] = make(map[reflect.Type]*dependencyNode)
			}

			context.DependencyGraph[inputType][returnType] = node
		}
	} else {
		// If there are no input types, then add the node to the start node list
		context.StartNodeList = append(context.StartNodeList, node)
	}
}

func (g *GimbapDependencyManager) ResolveDependencies(instanceMap map[reflect.Type]reflect.Value, providerList []*provider.Provider) {
	start := time.Now()

	// Create a new context
	context := NewGimbapDependencyManagerContext(instanceMap)

	// List the providers
	for _, p := range providerList {
		g.addProvider(context, p)
	}

	// If the starting node list is empty --> panic
	if len(context.StartNodeList) == 0 {
		g.logger.Panicf("Failed to resolve the dependencies. There is no starting node. At least one provider must not require any other provider")
	}

	// Instantiate the providers using the node data
	g.instantiateProviders(context)

	elapsed := time.Since(start)
	g.logger.Debugf("Dependency resolution took %s", elapsed)

	g.logger.Logf("Successfully resolved the dependencies for %d providers", len(providerList))
}

func (g *GimbapDependencyManager) OnStart() {
	g.logger.Log("Starting GimbapDependencyManager")
}

func (g *GimbapDependencyManager) OnStop() {
}

func NewGimbapDependencyManager() *GimbapDependencyManager {
	return &GimbapDependencyManager{
		logger: ecl.NewLogger(ecl.LoggerOption{
			Name: "GimbapDepManager",
		}),
	}
}
