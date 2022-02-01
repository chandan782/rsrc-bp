package main

import (
	"context"
	"fmt"
	"freecharge/rsrc-bp/api/helpers/logger"
	"freecharge/rsrc-bp/api/routes"
	"freecharge/rsrc-bp/configs"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// TODO: need to check for the cancelation of contexts
	// can be used to terminate the server using done
	pCtx := context.Background()
	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	config := configs.NewServerConfig(cancel)

	app := fiber.New(config.ServerConfig)

	// adding middleware
	app.Use(cors.New(config.CorsConfig))
	app.Use(helmet.New())

	// initialize router
	r, err := routes.NewRouteHandler(app, config)
	if err != nil {
		cancel()
	}

	r.NewRouter(app)

	// TODO: add request timeout and other configs to the server
	err = app.Listen(fmt.Sprintf(":%s", config.Port))
	if err != nil {
		cancel()
	}
	fmt.Println("successfully starting server")

	cChannel := make(chan os.Signal, 2)
	signal.Notify(cChannel, os.Interrupt, syscall.SIGTERM)

bLoop:
	for {
		select {
		case <-ctx.Done():
			// TODO: need to understand this more deaply.
			// list down the operations after Done is triggered

		case <-cChannel:
			logger.Warn("catch interrupted signal")
			// time.Sleep(sConfig.WaitTimeBeforeKill)
			break bLoop
		}
	}

	// TODO: add grpc shutting down handling
	err = app.Shutdown()
	if err != nil {
		// logging in server terminaltion
	}

	logger.Warn("shutting down the server :(")
}