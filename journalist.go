package main

import (
	"context"
	"log"
	"os"
  "fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"

	// "github.com/mrusme/journalist/ent"

	_ "github.com/mattn/go-sqlite3"
)

var fiberApp *fiber.App
var fiberLambda *fiberadapter.FiberLambda

func init() {
  log.Printf("Fiber cold start")
  fiberApp = fiber.New()

  fiberApp.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
  })

  fiberLambda = fiberadapter.New(fiberApp)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  return fiberLambda.ProxyWithContext(ctx, req)
}

func main() {
  appBindIp := os.Getenv("JOURNALIST_SERVER_BINDIP")
  appPort := os.Getenv("JOURNALIST_SERVER_PORT")
  functionName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")

  if functionName == "" {
    if appBindIp == "" {
      appBindIp = "127.0.0.1"
    }
    if appPort == "" {
      appPort = "8000"
    }
    log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%s", appBindIp, appPort)))
  } else {
    lambda.Start(Handler)
  }
}

