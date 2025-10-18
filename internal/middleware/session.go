package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/redis"
)

// function being used by csrf
func SessionStore() *session.Store {
	var sessionConfig = session.Config{
		Expiration:   24 * time.Hour,
		KeyLookup:    "cookie:session_id",
		KeyGenerator: utils.UUIDv4,
		Storage:      redis.New(), // using redis in order store session_ids
	}
	return session.New(sessionConfig)
}
