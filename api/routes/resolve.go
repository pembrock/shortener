package routes

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/pembrock/shortener/database"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if !errors.Is(err, redis.Nil) {
		//redirect to url
		return c.Redirect(value, 301)
	} else if errors.Is(err, redis.Nil) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "url not found"})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
