package apis

import (
	"net/http"

	"time"

	"pickrewardapi/base/ctx"
	"pickrewardapi/base/stats"

	"github.com/gin-gonic/gin"
)

func HandleWithQuery(group *gin.RouterGroup, httpMethod, relativePath, queryKeys string, handlers ...gin.HandlerFunc) {
	p := httpMethod + group.BasePath() + relativePath
	// Add pre-defined middlewared outside of custom ones
	handlers = append([]gin.HandlerFunc{
		// Stat is the very outer later track of all response code
		stats.Stat(p, queryKeys),
		// Have a recovery middleware inside to catch panic and log to gin
		// This also helps to convert panic to 500 return code so Stat() can track
		stats.Recovery(p),
	}, handlers...)

	group.Handle(httpMethod, relativePath, handlers...)
}

// Handle invokes RouterGroups's handle methid with counter middleware injected.
func Handle(group *gin.RouterGroup, httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	HandleWithQuery(group, httpMethod, relativePath, "", handlers...)
}

func BodySizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

func SetTimeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if d > time.Duration(0) {
			context := c.MustGet("ctx").(ctx.CTX)
			context, cancel := ctx.WithTimeout(context, d)
			defer cancel()
			c.Set("ctx", context)
		}
		c.Next()
	}
}
