package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"

	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/user"

	_ "github.com/mattn/go-sqlite3"
)

var fiberApp *fiber.App
var fiberLambda *fiberadapter.FiberLambda

func init() {
  log.Printf("Fiber cold start")
  fiberApp = fiber.New()

  entClient, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
  if err != nil {
    log.Fatalf("Failed initializing database: %v\n", err)
  }
  defer entClient.Close()
  if err := entClient.Schema.Create(context.Background()); err != nil {
    log.Fatalf("Failed initializing schema: %v\n", err)
  }

  var admin *ent.User
  admin, err = entClient.User.
    Query().
    Where(user.Username("admin")).
    Only(context.Background())
  if err != nil {
    appAdminPassword := os.Getenv("JOURNALIST_ADMIN_PASSWORD")
    if appAdminPassword == "" {
      appAdminPassword = "admin"
    }
    admin, err = entClient.User.Create().SetUsername("admin").SetPassword(appAdminPassword).Save(context.Background())
    if err != nil {
      log.Fatalf("Failed to query as well as create admin user: %v\n", err)
    }
  }

  log.Printf("Admin user: %s:%s\n", admin.Username, admin.Password)

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

