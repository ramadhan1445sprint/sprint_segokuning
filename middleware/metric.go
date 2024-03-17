package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/metric"
)

var sm *metric.ServerMetric

func init() {
	sm = metric.NewServerMetric()
}

func RecordDuration(ctx *fiber.Ctx) error {
	start := time.Now()

	err := ctx.Next()

	method := ctx.Method()
	path := ctx.Route().Path
	rawCode := ctx.Response().StatusCode()
	statusCode := strconv.Itoa(rawCode)

	if err != nil {
		if customError, ok := err.(customErr.CustomError); ok {
			statusCode = strconv.Itoa(customError.Status())
		} else if rawCode == fiber.StatusOK || rawCode == fiber.StatusCreated {
			statusCode = "500"
		}
	}

	elapsedDuration := time.Since(start).Seconds()

	fmt.Println(method, path, statusCode, elapsedDuration)
	sm.ReqDurationHist.WithLabelValues(method, path, statusCode).Observe(elapsedDuration)

	return err
}
