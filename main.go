package main

import (
	"ssego/messages"
	"ssego/registries"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
)

// user egistry - [global]
var user_registry *registries.UserRegistry = registries.NewUserRegitry()

func initSSE(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	//TODO: here extract the userid from JWT and create an entry in the user registry

	//keep alive msg repeater block
	for {
		select {
		case <-ticker.C:
			// Example: Emit an event every 5 seconds
			sendMsg := messages.TempMessage{
				Msg:  "Hello, this is a periodic message!",
				Type: "periodic",
			}
			if err := sendMsg.Emit(c); err != nil {
				return err
			}
		case <-c.Request().Context().Done():
			return nil
		}

	}
}

func main() {
	e := echo.New()

	//add - cron task to release the resgistri every day at 12pm [ midday ]
	c := cron.New()
	_, err := c.AddFunc("0 0 12 * * *", func() { // Cron expression for every day at 12:00 PM
		user_registry.Release()
	})
	if err != nil {
		e.Logger.Fatal(err)
	}
	c.Start()
	defer c.Stop()

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	e.GET("/startSSE", initSSE)
	e.Logger.Fatal(e.Start(":8080"))
}
