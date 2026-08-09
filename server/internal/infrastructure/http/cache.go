package http

import (
	"github.com/gin-gonic/gin"
)

// cacheControl returns a middleware that sets the Cache-Control response
// header for every request that reaches the wrapped handlers.
func cacheControl(header string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", header)
		c.Next()
	}
}

// pageCache is the caching policy for HTML documents: always revalidate,
// so updated pages are never served stale.
func pageCache() gin.HandlerFunc {
	return cacheControl("no-cache")
}
