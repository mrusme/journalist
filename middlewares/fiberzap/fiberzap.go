// Originally from https://gl.oddhunters.com/pub/fiberzap
// Copyright (apparently) by Ozgur Boru <boruozgur@yandex.com.tr>
// and "mert" (https://gl.oddhunters.com/mert)
package fiberzap

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Config defines the config for middleware
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Logger defines zap logger instance
	Logger *zap.Logger
}

// New creates a new middleware handler
func New(config Config) fiber.Handler {
	var (
		errPadding  = 15
		start, stop time.Time
		once        sync.Once
		errHandler  fiber.ErrorHandler
	)

	return func(c *fiber.Ctx) error {
		if config.Next != nil && config.Next(c) {
			return c.Next()
		}

		once.Do(func() {
			errHandler = c.App().Config().ErrorHandler
			stack := c.App().Stack()
			for m := range stack {
				for r := range stack[m] {
					if len(stack[m][r].Path) > errPadding {
						errPadding = len(stack[m][r].Path)
					}
				}
			}
		})

		start = time.Now()

		chainErr := c.Next()

		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		stop = time.Now()

		fields := []zap.Field{
			zap.Namespace("context"),
			zap.String("pid", strconv.Itoa(os.Getpid())),
			zap.String("time", stop.Sub(start).String()),
			zap.Object("response", Resp(c.Response())),
			zap.Object("request", Req(c)),
		}

		if u := c.Locals("userId"); u != nil {
			fields = append(fields, zap.Uint("userId", u.(uint)))
		}

		formatErr := ""
		if chainErr != nil {
			formatErr = chainErr.Error()
			fields = append(fields, zap.String("error", formatErr))
			config.Logger.With(fields...).Error(formatErr)

			return nil
		}

		config.Logger.With(fields...).Info("api.request")

		return nil
	}
}
