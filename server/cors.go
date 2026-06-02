package server

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func corsMiddleware() fiber.Handler {
	return cors.New(buildCorsConfig())
}

func buildCorsConfig() cors.Config {
	methods := []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"}
	headers := []string{"Origin", "Content-Type", "Accept", "Authorization"}

	origins := parseList(os.Getenv("CORS_ORIGINS"))
	allowVercelPreview := strings.EqualFold(os.Getenv("CORS_ALLOW_VERCEL_PREVIEW"), "true")

	if len(origins) == 0 && os.Getenv("GO_ENV") == "development" {
		return cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: methods,
			AllowHeaders: headers,
		}
	}

	if len(origins) == 0 {
		// Production fallback: tetap permissive agar tidak lock-out saat belum set env
		return cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: methods,
			AllowHeaders: headers,
		}
	}

	return cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			if origin == "" {
				return true
			}
			for _, allowed := range origins {
				if origin == allowed {
					return true
				}
			}
			if allowVercelPreview && strings.HasSuffix(origin, ".vercel.app") {
				return true
			}
			return false
		},
		AllowMethods: methods,
		AllowHeaders: headers,
	}
}

func parseList(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
