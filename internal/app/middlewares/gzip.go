package middlewares

import "github.com/labstack/echo/v4"

func Ommit(gzip echo.MiddlewareFunc, paths []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		gzipNext := gzip(next)
		return func(c echo.Context) error {
			// ignore prometeus
			path := c.Request().URL.Path

			for _, ommit := range paths {
				if path == ommit {
					return next(c)
				}
			}

			return gzipNext(c)
		}
	}
}
