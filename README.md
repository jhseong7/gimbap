<div align="center">
  <img src="docs/public/gimbap-logo.png" alt="GIMBAP Logo" />
</div>

# GIMBAP - Go Injection Management for Better Application Programming

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

**G**o
**I**njection
**M**anagement for
**B**etter
**A**pplication
**P**rogramming

## Why GIMBAP?

Gimbap is a Korean-style sea-weed roll with various ingredients insid. It is a simple and easy to make dish that is popular in Korea. The name GIMBAP is chosen as the acronym because it's one of my favourite dishes!.

Gimbap is a food where vegetables, fish and meat are rolled in a seaweed sheet with cooked rice. The seaweed sheet connects all the ingredients together and makes it easy to eat.

Like this GIMBAP aims to provide a simple solution to build web application in Go, connecting all components together in a simple and easy way.

## Disclaimer

GIMBAP is still in its **early stage of development**. It is not recommended to use this in production environment yet, and I am still working on to handle some of the key features that are missing or need improvement.

## Key Features

GIMBAP currently supports the following key features:

- Automatic Dependency Injection management
- API Endpoint management by code
- Flexible Server engine switching
- Containerized instance managing

### Automatic Dependency Injection Management

This is an example of adding a new dependency struct B to the constructor of A

> AS-IS

```go
package example

import "github.com/jhseong7/gimbap"

type A struct {}

func CreateA() *A {
  return &A{}
}

var ProviderA = gimbap.DefineProvider(
  gimbap.ProviderOption{Name: "A", Instantiator: CreateA},
)

var Module = gimbap.DefineModule(
  gimbap.ModuleOption{
    Name: "ExampleModule",
    Providers: []*gimbap.Provider{ProviderA}
  }
)
```

> TO-BE

```go
package example

import "github.com/jhseong7/gimbap"

type A struct {B *B}


func CreateA(B *B) *A {
  return &A{B: B}
}

var ProviderA = gimbap.DefineProvider(
  gimbap.ProviderOption{Name: "A", Instantiator: CreateA},
)

// ==== Add defs for B START =====
type B struct {}

func CreateB() *B {
  return &B{}
}

var ProviderB = gimbap.DefineProvider(
  gimbap.ProviderOption{Name: "B", Instantiator: CreateB},
)
// ==== Add defs for B END =====

var Module = gimbap.DefineModule(
  gimbap.ModuleOption{
    Name: "ExampleModule",
    Providers: []*gimbap.Provider{ProviderA, ProviderB} // Add the provider for B here
  }
)
```

Don't worry about changing the intialization process by hand for the new struct. GIMBAP will handle it for you.

### API Endpoint management by code

WIP

### Flexible Server engine switching

GIMBAP has an modularized server core that can be switched easily by preference

```go
import "github.com/jhseong7/gimbap"
import "github.com/jhseong7/gimbap/internal/engine/fiber_engine"
import "github.com/jhseong7/gimbap/internal/engine/echo_engine"

func main() {
  // Using GIN as http server
  a := gimbap.CreateApp(gimbap.AppOption{
    AppName:   "SampleApp",
		AppModule: SampleModule,
  })

  // Using fiber
  a = gimbap.CreateApp(gimbap.AppOption{
    AppName:   "SampleApp",
		AppModule: SampleModule,
    ServerEngine: fiber_engine.NewFiberHttpEngine()
  })

  // Using Echo
  a = gimbap.CreateApp(gimbap.AppOption{
    AppName:   "SampleApp",
		AppModule: SampleModule,
    ServerEngine: echo_engine.NewEchoHttpEngine()
  })
}
```

You can choose whatever web server you prefer as the core server easily.

## Documentation

Check for the documentation here.

[Documentation](https://go-gimbap.com)

## Related Projects

[Sample Repo](https://github.com/jhseong7/gimbap-sample)
[Config Provider](https://github.com/jhseong7/gimbap-config)

## Goals

This framework provides the following features in the future

- Interceptor, Guard support
- Template engine support
- More to be added...

## Third-Party Libraries

This project uses the following third-party libraries:

- **Library Name:** uber-go/fx
  - **Purpose:** Used managing dependency management in the `FxDependencyManager`
  - **License:** MIT License. [Link](https://github.com/uber-go/fx/blob/master/LICENSE)
  - **Link:** [https://github.com/uber-go/fx](https://github.com/uber-go/fx)
- **Library Name:** gin-gonic/gin
  - **Purpose:** Used the internal server engine `GinHttpEngine`
  - **License:** MIT License. [Link](https://github.com/gin-gonic/gin/blob/master/LICENSE)
  - **Link:** [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- **Library Name:** gofiber/fiber
  - **Purpose:** Used the internal server engine `FiberHttpEngine`
  - **License:** MIT License. [Link](https://github.com/gofiber/fiber/blob/master/LICENSE)
  - **Link:** [https://github.com/gofiber](https://github.com/gofiber)
- **Library Name:** labstack/echo
  - **Purpose:** Used the internal server engine `EchoHttpEngine`
  - **License:** MIT License. [Link](https://github.com/labstack/echo/blob/master/LICENSE)
  - **Link:** [https://github.com/labstack/echo](https://github.com/labstack/echo)

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/jhseong7/gimbap/tree/main/LICENSE) file for details.
