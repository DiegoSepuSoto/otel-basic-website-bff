package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/DiegoSepuSoto/basic-website-bff/src/tracing"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
)

type OrderAPIResponse struct {
	OrderID     string `json:"orderID"`
	OrderStatus string `json:"orderStatus"`
	Customer    struct {
		Name       string `json:"name"`
		LastName   string `json:"lastName"`
		CustomerID string `json:"customerID"`
	} `json:"customer"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	ctx := context.Background()

	tp, err := tracing.InitTelemetryExporter(ctx)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("error shutting down tracer provider: %v", err)
		}
	}()

	r := echo.New()
	r.HideBanner = true
	r.Use(otelecho.Middleware("basic-website-bff"))

	r.GET("/order", getOrderHandler)

	_ = r.Start(":8082")
}

func getOrderHandler(c echo.Context) error {
	ctx, span := otel.Tracer(tracing.TracerName).Start(c.Request().Context(), "getOrderHandler")
	defer span.End()

	time.Sleep(500 * time.Millisecond)

	orderID := c.QueryParam("orderID")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/order?orderID=%s",
			os.Getenv("ORDER_API_HOST"),
			orderID,
		),
		nil,
	)
	if err != nil {
		log.Errorf("error making API request %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error making API request")
	}

	response, err := tracing.HTTPClient.Do(req)
	if err != nil {
		log.Errorf("error doing API request %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error doing API request")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("error reading API response %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error reading API response")
	}

	var orderAPIResponse OrderAPIResponse
	err = json.Unmarshal(body, &orderAPIResponse)
	if err != nil {
		log.Errorf("error unmarshalling API response %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error unmarshalling API response")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("La orden de %s %s está %s actualmente", orderAPIResponse.Customer.Name,
			orderAPIResponse.Customer.LastName, orderAPIResponse.OrderStatus),
	})
}
