// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package limiter

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber"
)

// Config ...
type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool
	// Timeout in seconds on how long to keep records of requests in memory
	// Default: 60
	Timeout int
	// Max number of recent connections during `Timeout` seconds before sending a 429 response
	// Default: 10
	Max int
	// Message
	// default: "Too many requests, please try again later."
	Message string
	// StatusCode
	// Default: 429 Too Many Requests
	StatusCode int
	// Key allows to use a custom handler to create custom keys
	// Default: func(c *fiber.Ctx) string {
	//   return c.IP()
	// }
	Key func(*fiber.Ctx) string
	// Handler is called when a request hits the limit
	// Default: func(c *fiber.Ctx) {
	//   c.Status(cfg.StatusCode).SendString(cfg.Message)
	// }
	Handler func(*fiber.Ctx)
}

// New ...
func New(config ...Config) func(*fiber.Ctx) {
	// mutex for parallel read and write access
	mux := &sync.Mutex{}
	// Init config
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 60
	}
	if cfg.Max == 0 {
		cfg.Max = 10
	}
	if cfg.Message == "" {
		cfg.Message = "Too many requests, please try again later."
	}
	if cfg.StatusCode == 0 {
		cfg.StatusCode = 429
	}
	if cfg.Key == nil {
		cfg.Key = func(c *fiber.Ctx) string {
			return c.IP()
		}
	}
	if cfg.Handler == nil {
		cfg.Handler = func(c *fiber.Ctx) {
			c.Status(cfg.StatusCode).SendString(cfg.Message)
		}
	}
	// Limiter settings
	var hits = make(map[string]int)
	var reset = make(map[string]int)
	var timestamp = int(time.Now().Unix())
	// Update timestamp every second
	go func() {
		for {
			timestamp = int(time.Now().Unix())
			time.Sleep(1 * time.Second)
		}
	}()
	// Reset hits every cfg.Timeout
	go func() {
		for {
			// For every key in reset
			for key := range reset {
				// If resetTime exist and current time is equal or bigger
				if reset[key] != 0 && timestamp >= reset[key] {
					// Reset hits and resetTime
					mux.Lock()
					hits[key] = 0
					reset[key] = 0
					mux.Unlock()
				}
			}
			// Wait cfg.Timeout
			time.Sleep(time.Duration(cfg.Timeout) * time.Second)
		}
	}()
	return func(c *fiber.Ctx) {
		// Filter request to skip middleware
		if cfg.Filter != nil && cfg.Filter(c) {
			c.Next()
			return
		}
		// Get key (default is the remote IP)
		key := cfg.Key(c)
		mux.Lock()
		// Increment key hits
		hits[key]++
		// Set unix timestamp if not exist
		if reset[key] == 0 {
			reset[key] = timestamp + cfg.Timeout
		}
		// Get current hits
		hitCount := hits[key]
		// Calculate when it resets in seconds
		resetTime := reset[key] - timestamp
		mux.Unlock()
		// Set how many hits we have left
		remaining := cfg.Max - hitCount
		// Check if hits exceed the cfg.Max
		if remaining < 0 {
			// Call Handler func
			cfg.Handler(c)
			// Return response with Retry-After header
			// https://tools.ietf.org/html/rfc6584
			c.Set("Retry-After", strconv.Itoa(resetTime))
			return
		}
		// We can continue, update RateLimit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(cfg.Max))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Set("X-RateLimit-Reset", strconv.Itoa(resetTime))
		// Bye!
		c.Next()
	}
}
