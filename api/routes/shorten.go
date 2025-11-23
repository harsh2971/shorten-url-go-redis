package routes

import (
	"os"
	"shorten-url-go-redis/database"
	"shorten-url-go-redis/helpers"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"custom_short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// ShortenURL is a function that shortens a URL
// It takes a request body and returns a response body
// The request body is a JSON object with the following fields:
// - URL: the URL to shorten
// - CustomShort: the custom short URL
// - Expiry: the expiry time of the URL
// The response body is a JSON object with the following fields:
// - URL: the URL to shorten
// - CustomShort: the custom short URL
// - Expiry: the expiry time of the URL
func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "cannot parse JSON",
		})
	}

	//implement rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()

	// check if the IP is in the database by key (IP)
	// if not, set the IP with the value of API_QUOTA for 30 minutes
	// if it is, check if the rate limit is exceeded or value is less than 1
	// if it is, return service unavailable with time remaining in minutes to reset the rate limit
	value, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second)
	} else {
		valInt, _ := strconv.Atoi(value)
		if valInt < 1 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            true,
				"message":          "rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}

	}

	//check if URL is valid
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid URL",
		})
	}

	//check for domain error - remove http:// or https:// from the URL and check if it's valid
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   true,
			"message": "domain error",
		})
	}

	//enforce https, SSL/TLS
	body.URL = helpers.EnforceHTTP(body.URL)

	//generate custom short
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	// Validate custom short if provided
	if body.CustomShort != "" {
		if len(body.CustomShort) < 6 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "custom short is too short",
			})
		} else if len(body.CustomShort) > 15 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "custom short is too long",
			})
		}
		// Check if custom short already exists
		exist, _ := r.Get(database.Ctx, id).Result()
		if exist != "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "custom short already exists",
			})
		}
	} else {
		// For auto-generated IDs, check if it already exists and regenerate if needed
		exist, _ := r.Get(database.Ctx, id).Result()
		for exist != "" {
			id = uuid.New().String()[:6]
			exist, _ = r.Get(database.Ctx, id).Result()
		}
	}

	// Set expiry default
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	// Save to Redis
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error saving URL to database",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  20,
		XRateLimitReset: 30 * time.Minute,
	}
	//decrement the rate limit
	r2.Decr(database.Ctx, c.IP())

	val, _ := r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	return c.Status(fiber.StatusOK).JSON(resp)
}
