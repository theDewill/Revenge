package Revenge

import (
	"fmt"
	"net/http"
	"ssego/messages"
	"ssego/registries"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
)

// RV injector to the SSEhandler
func GenerateCustomDataMiddleware(RV *RevengeRoot) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Store it in the context
			c.Set("RVOBJ", RV)

			//TODO: this is temp and remove this
			c.Set("ugrpid", 1)
			c.Set("uid", 1)
			// Call the next handler in the chain
			return next(c)
		}
	}
}

// user egistry - [global]

type Server struct {
	port string
	Res  *echo.Echo
}

type Runner struct{}

type Portion struct {
}

//TODO: this will be added to the revenge once server pool and load balncer is placed

type RevengeRoot struct {
	server        Server
	task_manger   *cron.Cron
	user_registry *registries.UserRegistry
}

func New(port string) *RevengeRoot {
	return &RevengeRoot{
		user_registry: registries.NewUserRegitry(),
		server: Server{
			Res:  echo.New(),
			port: port,
		},
		task_manger: cron.New(),
	}
}

func (RV *RevengeRoot) Commence() {

	_, err := RV.task_manger.AddFunc("0 12 * * *", func() { // Cron expression for every day at 12:00 PM
		RV.user_registry.Release()
	})
	if err != nil {
		RV.server.Res.Logger.Fatal(err)
	}
	RV.task_manger.Start()
	defer RV.task_manger.Stop()

	RV.server.Res.Use(GenerateCustomDataMiddleware(RV))

	RV.server.Res.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	//Rest Point
	RV.server.Res.GET("/stconn", SSEhandler)

	RV.server.Res.GET("/testUpdate", test_handler)

	RV.server.Res.Logger.Fatal(RV.server.Res.Start(":" + RV.server.port))
}

func Web_Socket() interface{} {

}

func test_handler(c echo.Context) error {
	//performing db operations - chanign the state of DB
	new_entry := "Hello World" // new entry
	UR := c.Get("RVOBJ").(RevengeRoot).user_registry

	msg := &messages.UpdateMessage{
		Msg:        "Hello, this is a periodic test msg!",
		Type:       "periodic",
		Data:       "new entry",
		Event:      "new_entry",
		Action:     "insert",
		Dependents: "none",
		Component:  "none",
	}

	UR.SendUpdates(1, 1, msg)

	return c.JSON(http.StatusOK, new_entry)
}

func SSEhandler(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("SSE handler called")

	//TODO: here extract the userid from JWT and create an entry in the user registry]

	var RV *RevengeRoot = c.Get("RVOBJ").(*RevengeRoot)

	tmp_user := RV.user_registry.LoadUser(c.Get("ugrpid").(int), c.Get("uid").(int))

	fmt.Printf("passsed loading phase")

	//keep alive msg repeater block
	for {
		select {

		//SSE activator function
		case message := <-tmp_user.Sse_channel:
			fmt.Println("Message Received for other users")
			c.Response().WriteHeader(http.StatusOK)
			if _, err := c.Response().Write([]byte(message.(string))); err != nil {
				return err
			}
			c.Response().Flush()

		case <-ticker.C:
			// Example: Emit an event every 5 seconds
			fmt.Println("Message sent for temp")
			sendMsg := messages.TempMessage{
				Msg:  "Hello, this is a periodic message!",
				Type: "periodic",
			}
			if err := sendMsg.EmitTmp(c); err != nil {
				return err
			}
		case <-c.Request().Context().Done():
			return nil
		}

	}
}

//OLD-Main

func main_old() {
	e := echo.New()

	//GQL
	fields := graphql.Fields{
		"startSSE": &graphql.Field{
			Type:        graphql.String,
			Description: "Start an SSE stream",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Simulate the SSE initialization logic
				return "SSE Initialized", nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, gerr := graphql.NewSchema(schemaConfig)
	if gerr != nil {
		e.Logger.Fatal(gerr)
	}

	//add - cron task to release the resgistri every day at 12pm [ midday ]
	// c := cron.New()
	// _, err := c.AddFunc("0 0 12 * * *", func() { // Cron expression for every day at 12:00 PM
	// 	user_registry.Release()
	// })
	// if err != nil {
	// 	e.Logger.Fatal(err)
	// }
	// c.Start()
	// defer c.Stop()

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	//G-Point
	e.POST("/graphql", func(c echo.Context) error {
		var p graphql.Params
		if err := c.Bind(&p); err != nil {
			return err
		}
		p.Schema = schema
		result := graphql.Do(p)
		if len(result.Errors) > 0 {
			return c.JSON(400, result.Errors)
		}
		return c.JSON(200, result)
	})

	//Rest Point
	e.GET("/startSSE", SSEhandler)
	e.Logger.Fatal(e.Start(":8080"))
}
