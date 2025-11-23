package routes

import (
	"github.com/gofiber/fiber/v2"
	"shorten-url-go-redis/database"
	"github.com/redis/go-redis/v9"
)

// ResolveURL is a function that resolves a URL
// It takes a URL parameter and returns a redirect response
// The URL parameter is the short URL to resolve
// The response is a redirect to the original URL
// If the URL is not found, it returns a 404 error
// If the URL is found, it returns a redirect to the original URL
// If the URL is not found, it returns a 404 error
func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	r:=database.CreateClient(0)
	defer r.Close()

	// check if the URL is in the database by key (url)
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"message": "short not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message": "error getting value from database",
		})
	}

	rInr:=database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")
	return c.Redirect(value, 302)


}