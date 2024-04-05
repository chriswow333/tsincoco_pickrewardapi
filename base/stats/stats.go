package stats

import (
	"net/http/httputil"

	"pickrewardapi/base/ctx"
	"pickrewardapi/base/goroutine"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Stat(staticPath string, queryKeys string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sets static path to be logged in Log() middlerware
		c.Set("staticPath", staticPath)

		// Put path into metric tags.
		// NOTE be cautious to add tags because it could greatly
		// increase the number of custom metrics using this tags.

		//	tags := []string{"path", staticPath}
		//

	}
}

func Recovery(staticPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				stack := goroutine.Stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				c.MustGet("ctx").(ctx.CTX).WithFields(logrus.Fields{
					"path":    staticPath,
					"request": string(httprequest),
					"err":     err,
					"stack":   string(stack),
				}).Info("panic")

				c.AbortWithStatus(500)
			}
		}()
	}
}
