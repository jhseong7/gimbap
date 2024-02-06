# NASSI Golang

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

(N)estJS
(A)nd
(S)pring
(S)imilarly
(I)n
Golang

## Introduction

This is a Go Web framework inspired by Nest.JS and Spring, introducing similar patterns of the frameworks and adapt it to the Golang style pattern

Even major web frameworks like `gin` intializes routing handlers as below

```golang
r := gin.Default()

r.GET("/test", a.Handler1)
r.POST("/test", a,Handler2)
```

which is similar to the classic Node.JS's `express`.

```js
const app = express();
const a = new A();

app.get("/test", a.handler1);
app.post("/test", a.handler2);
```

The huge advantages of frameworks like NestJS and Spring, are its AOP patterns and strict model hierarchy that allows simpler design patterns.

For example, the simple code above to define router handlers would be like below

```ts
@Controller()
class TestController {
  @Get("test")
  public async handler1() {
    /* ... */
  }

  @Post("test")
  public async handler1() {
    /* ... */
  }
}
```

```java
@RestController
class TestController {
  @GetMapping("/test")
  public String handler1() {/* ... */}

  @PostMapping("/test")
  public String handler2() {/* ... */}
}
```

For small applications, the current pattern should work fine, but as projects grow and get bigger, AOPs help to manage the code in a managed way.

## Goals

This framework provides the following features:

- Module based Dependency Injection (DI)
  - This will use `uber/fx` for runtime DI
- `Controller` adapter pattern.
  - Adapters can be replaced but the framework will provide a constant handler spec
  - GIN + a
- `Injectable` interface for module based DI
  - All services, repos that require DI must be specified in the right way
- ORM adapters
  - external module. This will be vendor-locked to that specific module.
  - initial support will start with `gORM`
- Config manager (dotenv, config service)
- Container based server control
  - Container isolation like NestJS or Spring
  - May be a bean-like config management
- Panic handle middleware support
- Interceptor, Guard support
