package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	// "github.com/gofiber/template/html"

	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/user"
	"github.com/mrusme/journalist/journalistd"

	"github.com/mrusme/journalist/api"
	"github.com/mrusme/journalist/web"

	_ "github.com/mattn/go-sqlite3"
  _ "github.com/lib/pq"
  _ "github.com/go-sql-driver/mysql"
)

//go:embed views/*
var viewsfs embed.FS

var fiberApp *fiber.App
var fiberLambda *fiberadapter.FiberLambda
var entClient *ent.Client

func init() {
  log.Printf("Fiber cold start")

  fiberLambda = fiberadapter.New(fiberApp)
}

func Handler(
  ctx context.Context,
  req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
  return fiberLambda.ProxyWithContext(ctx, req)
}

func main() {
  var err error
  defer entClient.Close()

  config, err := journalistd.Cfg()
  if err != nil {
    log.Panic(err.Error())
  }

  engine := web.NewFileSystem(http.FS(viewsfs), ".html")
  fiberApp = fiber.New(fiber.Config{
    Views: engine,
  })

  entClient, err = ent.Open(config.Database.Type, config.Database.Connection)
  if err != nil {
    log.Fatalf("Failed initializing database: %v\n", err)
  }
  if err := entClient.Schema.Create(context.Background()); err != nil {
    log.Fatalf("Failed initializing schema: %v\n", err)
  }

  var admin *ent.User
  admin, err = entClient.User.
    Query().
    Where(user.Username(config.Admin.Username)).
    Only(context.Background())
  if err != nil {
    admin, err = entClient.User.
      Create().
      SetUsername(config.Admin.Username).
      SetPassword(config.Admin.Password).
      SetRole("admin").
      Save(context.Background())
    if err != nil {
      log.Fatalf("Failed to query as well as create admin user: %v\n", err)
    }
  }

  if admin.Password == "admin" {
    log.Printf("Admin user: %s:%s\n", admin.Username, admin.Password)
  } else {
    log.Printf("Admin user: %s:xxxxxxxx\n", admin.Username)
  }

  fiberApp.Use(logger.New())
  api.Register(&config, fiberApp, entClient)
  web.Register(&config, fiberApp, entClient)


  functionName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")

  if functionName == "" {
    go func() {
      jd := journalistd.New(&config, entClient)

      // for {
        time.Sleep(time.Second * 30)
        errs := jd.RefreshAll()
        log.Printf("\n%v\n", errs)
        time.Sleep(time.Second * 30)
      // }
    }()

    log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%s", config.Server.BindIP, config.Server.Port)))
  } else {
    lambda.Start(Handler)
  }
}

