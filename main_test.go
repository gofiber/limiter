// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package limiter

import (
	"github.com/gofiber/fiber"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func Test_Concurrency(t *testing.T) {
	app := fiber.New()
	app.Use(New(Config{Max: 100}))
	app.Get("/", func(ctx *fiber.Ctx) {
		// random delay between the requests
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		ctx.Send("Hello tester!")
	})

	var wg sync.WaitGroup
	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Unexpected status code %v", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil || "Hello tester!" != string(body) {
				t.Fatalf("Unexpected body %v", string(body))
			}
		}()
	}

	wg.Wait()
}
