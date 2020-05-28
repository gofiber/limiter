# Limiter

![Release](https://img.shields.io/github/release/gofiber/limiter.svg)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)
![Test](https://github.com/gofiber/limiter/workflows/Test/badge.svg)
![Security](https://github.com/gofiber/limiter/workflows/Security/badge.svg)
![Linter](https://github.com/gofiber/limiter/workflows/Linter/badge.svg)

### Install
```
go get -u github.com/gofiber/fiber
go get -u github.com/gofiber/limiter
```
### Example
```go
package main

import (
  "github.com/gofiber/fiber"
  "github.com/gofiber/limiter"
)

func main() {
  app := fiber.New()

  // 3 requests per 10 seconds max
  cfg := limiter.Config{
    Timeout: 10,
    Max: 3,
  }

  app.Use(limiter.New(cfg))

  app.Get("/", func(c *fiber.Ctx) {
    c.Send("Welcome!")
  })

  app.Listen(3000)
}
```
### Test
```curl
curl http://localhost:3000
curl http://localhost:3000
curl http://localhost:3000
curl http://localhost:3000
```
### Third party implementations
[Limiter with Redis suppport](https://github.com/Shareed2k/fiber_limiter)
