package apiserver

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type Stats struct {
}

func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		log.Info(c.Path())
		return nil
	}
}
