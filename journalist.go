package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"

	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/user"
	"github.com/mrusme/journalist/journalistd"

	"github.com/mrusme/journalist/api"
	"github.com/mrusme/journalist/middlewares/fiberzap"
	"github.com/mrusme/journalist/web"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed views/*
var viewsfs embed.FS

//go:embed favicon.ico
var favicon embed.FS

var fiberApp *fiber.App
var fiberLambda *fiberadapter.FiberLambda
var entClient *ent.Client

func init() {
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
    panic(err)
  }

  logger, _ := zap.NewProduction()
  defer logger.Sync()
  // TODO: Use sugarLogger
  // sugar := logger.Sugar()

  engine := web.NewFileSystem(http.FS(viewsfs), ".html")
  fiberApp = fiber.New(fiber.Config{
    Views: engine,
  })
  fiberApp.Use(fiberzap.New(fiberzap.Config{
    Logger: logger,
  }))


  entClient, err = ent.Open(config.Database.Type, config.Database.Connection)
  if err != nil {
    logger.Error(
      "Failed initializing database",
      zap.Error(err),
    )
  }
  if err := entClient.Schema.Create(context.Background()); err != nil {
    logger.Error(
      "Failed initializing schema",
      zap.Error(err),
    )
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
      logger.Error(
        "Failed query/create admin user",
        zap.Error(err),
      )
    }
  }

  if admin.Password == "admin" {
    logger.Debug(
      "Admin user",
      zap.String("username", admin.Username),
      zap.String("password", admin.Password),
    )
  } else {
    logger.Debug(
      "Admin user",
      zap.String("username", admin.Username),
      zap.String("password", "xxxxxx"),
    )
  }

  api.Register(&config, fiberApp, entClient, logger)
  web.Register(&config, fiberApp, entClient, logger)

  fiberApp.Get("/favicon.ico", func(ctx *fiber.Ctx) error {
    fi, err := favicon.Open("favicon.ico")
    if err != nil {
      return ctx.SendStatus(fiber.StatusInternalServerError)
    }
    return ctx.SendStream(fi)
  })

  functionName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")

  if config.Feeds.AutoRefresh != "" {
    interval, err := strconv.Atoi(config.Feeds.AutoRefresh)
    if err != nil {
      logger.Fatal(
        "Feeds.AutoRefresh is not a valid number (seconds)",
        zap.Error(err),
      )
    }

    if functionName == "" {
      go func() {
        jd := journalistd.New(&config, entClient)

        time.Sleep(time.Second * 10)
        for {
          logger.Debug(
            "Running RefreshAll to refresh all feeds",
          )
          errs := jd.RefreshAll()
          if len(errs) > 0 {
            logger.Error(
              "RefreshAll completed with errors",
              zap.Errors("errors", errs),
            )
          }
          time.Sleep(time.Second * time.Duration(interval))
        }
      }()
    } else {
      logger.Warn(
        "Journalist won't start the feed auto refresh thread " +
        "while it is running as a Lambda function",
      )
    }
  }

  if functionName == "" {
    listenAddr := fmt.Sprintf(
      "%s:%s",
      config.Server.BindIP,
      config.Server.Port,
    )
    logger.Fatal(
      "Server failed",
      zap.Error(fiberApp.Listen(listenAddr)),
    )
  } else {
    lambda.Start(Handler)
  }
}

