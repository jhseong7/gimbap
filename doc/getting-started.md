[Return to Documentation](./documentation.md)

# Getting Started

GIMBAP is a framework aimed to quickly build and manage web applications in Go. It borrows some concepts from Spring Boot and NestJS and adapts them into a pattern that is more suitable for Go.

The following guide will help you quickly get started with the concepts of GIMBAP and make it possible to build a simple module managed web application.

## Installation

You can either implicitly install by importing the package like this

```golang
import "github.com/jhseong7/gimbap"
```

or explicitly install the package yourself

```shell
go get -u github.com/jhseong7/gimbap
```

## Initial Folder setup

Since GIMBAP is a framework that is designed to manage code in a modularized way, it is recommended to use a folder structure that is similar to frameworks of other languages.

The options recommended are as follows:

1. sub-folder by features
   - e.g. user, post, comment, etc.
   - all the code related to the feature should be in the same folder
2. sub-folder by types (models, services, controllers, etc.)
   - e.g. models, services, controllers, etc.
   - all the code related to the type should be in the same folder

However any folder structure you prefer is fine as long as you can manage the code well. 😃

### Sample Folder Structures

#### By Features

> main.go in the root folder

```shell
GIMBAP-project/
├── user/
│   ├── module.go
│   ├── user-service.go
│   ├── user-controller.go
│   └── user-repository.go
├── post/
│   ├── module.go
│   ├── post-service.go
│   ├── post-controller.go
│   └── post-repository.go
├── comment/
│   ├── module.go
│   ├── comment-service.go
│   ├── comment-controller.go
│   └── comment-repository.go
├── go.mod
├── go.sum
└── main.go
```

> main.go in the cmd folder. The modules can stay in the root or moved to internal/ folder depending on the preference.

```shell
GIMBAP-project/
├── cmd/
│   └── app/
│       └── main.go
├── app/
│   ├── module.go
│   ├── user/
│   │   ├── module.go
│   │   ├── user-service.go
│   │   ├── user-controller.go
│   │   └── user-repository.go
│   ├── post/
│   │   ├── module.go
│   │   ├── post-service.go
│   │   ├── post-controller.go
│   │   └── post-repository.go
│   └── comment/
│       ├── module.go
│       ├── comment-service.go
│       ├── comment-controller.go
│       └── comment-repository.go
├── go.mod
└── go.sum
```

#### By Types

```shell
GIMBAP-project/
├── controller/
│   ├── module.go
│   ├── user-controller.go
│   ├── post-controller.go
│   └── comment-controller.go
├── service/
│   ├── module.go
│   ├── user-service.go
│   ├── post-service.go
│   └── comment-service.go
├── repository/
│   ├── module.go
│   ├── user-repository.go
│   ├── post-repository.go
│   └── comment-repository.go
├── go.mod
├── go.sum
└── main.go
```

> main.go in the cmd folder. The modules can stay in the root or moved to internal/ folder depending on the preference.

```shell
GIMBAP-project/
├── cmd/
│   └── app/
│       └── main.go
├── app/
│   ├── module.go
│   ├── controller/
│   │   ├── module.go
│   │   ├── user-controller.go
│   │   ├── post-controller.go
│   │   └── comment-controller.go
│   ├── service/
│   │   ├── module.go
│   │   ├── user-service.go
│   │   ├── post-service.go
│   │   └── comment-service.go
│   └── repository/
│       ├── module.go
│       ├── user-repository.go
│       ├── post-repository.go
│       └── comment-repository.go
├── go.mod
└── go.sum
```

## Module setup

It is recommended to use a root module which combines all the module