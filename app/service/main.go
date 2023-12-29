package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/jackc/pgx/v5"
	"github.com/julioc98/ledger/app/service/api"
	v1 "github.com/julioc98/ledger/app/service/api/v1"
	"github.com/julioc98/ledger/domain/account"
	"github.com/julioc98/ledger/gateway/pg"
	"github.com/julioc98/ledger/gateway/redisx"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/julioc98/ledger/app/service/docs"
)

type config struct {
	ServerAddr         string        `conf:"env:SERVER_ADDR,default:0.0.0.0:3000"`
	ServerReadTimeout  time.Duration `conf:"env:SERVER_READ_TIMEOUT,default:30s"`
	ServerWriteTimeout time.Duration `conf:"env:SERVER_WRITE_TIMEOUT,default:30s"`
	DataBaseURL        string        `conf:"env:DATABASE_URL,default:postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	RedisURL           string        `conf:"env:REDIS_URL,default:redis://:@localhost:6379/0"`
}

// @title Ledger API
// @version 1.0
// @description API for managing ledger accounts
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
func main() {
	var prefix = os.Getenv("ENV_PREFIX")
	var cfg config
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		}
		fmt.Printf("parsing config: %s \n", err.Error())
		return
	}

	conn, err := pgx.Connect(context.Background(), cfg.DataBaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	log.Printf("Connected to Redis: %s", cfg.RedisURL)
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to redis: %v\n", err)
		return
	}

	client := redis.NewClient(opt)

	status := client.Ping(context.Background())
	if status.Err() != nil {
		fmt.Fprintf(os.Stderr, ":/ Unable to connect to redis: %v\n ", err)
		return
	}

	// Repository
	accountRepo := pg.NewAccountPgxRepository(conn)
	accountCacheRepo := redisx.NewAccountRedisRepositoryDecorator(client, accountRepo)

	// Use cases
	depositUC := account.NewDepositUseCase(accountRepo)
	transfUC := account.NewTransferUseCase(accountRepo, accountRepo)
	balanceUC := account.NewBalanceUseCase(accountRepo)
	transfHisUC := account.NewTransfersHistoryUseCase(accountCacheRepo)

	// Handlers
	apiv1 := v1.API{
		CreateDepositHandler:       v1.CreateDepositHandler(depositUC),
		CreateTransferHandler:      v1.CreateTransferHandler(transfUC),
		GetBalanceHandler:          v1.GetBalanceHandler(balanceUC),
		GetTransfersHistoryHandler: v1.GetTransfersHistoryHandler(transfHisUC),
	}

	// Init server
	router := api.NewServer()

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
	))

	// Routing
	apiv1.Routes(router)

	// Server
	server := http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      router,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}

	log.Printf("Server listening on %s", cfg.ServerAddr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
