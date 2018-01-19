package main

import (
	"fmt"
	"log"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

func main() {
	app := iris.New()
	// Load all templates from the "./views" folder
	// where extension is ".html" and parse them
	// using the standard `html/template` package.
	app.RegisterView(iris.HTML("./views", ".html"))

	// Use Middleware
	app.Use(func(ctx iris.Context) {
		log.Println("Simple Middlewares.")
		ctx.Next()
	})

	// Method:    GET
	// Resource:  http://localhost:8080
	app.Get("/", func(ctx iris.Context) {
		// Bind: {{.message}} with "Hello world!"
		ctx.ViewData("message", "Hello world!")
		// Render template file: ./views/hello.html
		ctx.View("hello.html")
	})

	app.Get("/socket", func(ctx iris.Context) {
		ctx.View("websocket.html")
	})

	// Method:    GET
	// Resource:  http://localhost:8080/user/42
	//
	// Need to use a custom regexp instead?
	// Easy,
	// just mark the parameter's type to 'string'
	// which accepts anything and make use of
	// its `regexp` macro function, i.e:
	// app.Get("/user/{id:string regexp(^[0-9]+$)}")
	app.Get("/user/{id:long}", func(ctx iris.Context) {
		userID, _ := ctx.Params().GetInt64("id")
		ctx.Writef("User ID: %d", userID)
	})

	app.Get("/json/{str:string}", func(ctx iris.Context) {
		str := ctx.Params().Get("str")
		ctx.JSON(iris.Map{
			"msg": fmt.Sprintf("This is GET [string=%s]", str),
		})
	})

	app.Post("/json", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"msg": "This is Post",
		})
	})

	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	ws.OnConnection(handleConnection)

	app.Get("/sync", authHandler, ws.Handler())

	app.Any("/iris-ws.js", websocket.ClientHandler())

	// Start the server using a network address.
	app.Run(iris.Addr(":8080"))
}

func authHandler(ctx iris.Context) {
	token := ctx.URLParam("token")
	log.Println("Auth: " + token)
	isAuth := (token == "myToken")
	if isAuth {
		ctx.Next()
	}
	ctx.EndRequest()
}

func handleConnection(c websocket.Connection) {
	log.Println("Status: Connected")

	// Use default WebSocket API Javascript
	// c.OnMessage(func(b []byte) {
	// 	msg := string(b)

	// 	fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)

	// 	c.To(websocket.Broadcast).EmitMessage(b)
	// })

	c.On("echo", func(msg string) {
		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)

		c.To(websocket.All).Emit("echo", "From server -> "+msg)
	})

	c.OnDisconnect(func() {
		log.Println("Status: Disconnected")
	})
}
