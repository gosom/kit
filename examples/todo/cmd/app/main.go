package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gosom/kit/es"
	"github.com/gosom/kit/es/eshttp"
	"github.com/gosom/kit/es/postgres"
	"github.com/gosom/kit/examples/todo"
	"github.com/gosom/kit/examples/todo/api"
	"github.com/gosom/kit/examples/todo/assets"
	"github.com/gosom/kit/logging"
	"github.com/gosom/kit/sqldb"
	"github.com/gosom/kit/web"
)

func main() {
	logger := logging.New("zerolog", logging.DEBUG, os.Stderr)
	logging.SetDefault(logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		cancel()
	}()
	if err := run(ctx); err != nil {
		logging.Error("error in run methdo", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	registry := es.NewRegistry()
	todo.Register(registry)

	db, err := getDb("postgres://todo:todo@localhost:5432/todo?sslmode=disable")
	if err != nil {
		return err
	}

	store := postgres.NewEventStore(db)
	if err := store.Migrate(ctx); err != nil {
		return err
	}

	if err := sqldb.Migrate(ctx, db, "", assets.Migrations); err != nil {
		return err
	}

	commandProcessor, err := es.NewCommandProcessor(
		2,
		store,
		registry,
		todo.DOMAIN,
	)
	if err != nil {
		return err
	}

	dispatcher := postgres.NewCommandDispatcher(todo.DOMAIN, store)
	webServer := getWebServer(store, registry, dispatcher)

	projectionBuilder := todo.NewProjectionBuilder(db, registry)

	appSvc, err := es.New(
		es.WithLogger(logging.Get().Level(logging.DEBUG)),
		es.WithEventStore(store),
		es.WithCommandProcessor(commandProcessor),
		es.WithWebServer(webServer),
		es.WithPublishers(projectionBuilder),
	)
	if err != nil {
		return err
	}

	return appSvc.Start(ctx)

}

func getDb(dsn string) (*sqldb.DB, error) {
	dbconn := sqldb.NewDB("postgres", dsn)
	return dbconn, dbconn.Open()
}

func getWebServer(store es.EventStore, registry *es.Registry, dispatcher es.CommandDispatcher) *web.HttpServer {
	routerCfg := web.RouterConfig{}
	mux := web.NewRouter(routerCfg)
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		web.JSON(w, r, http.StatusOK, map[string]string{"status": "ok"})
	})

	api.RegisterHandlers(mux, dispatcher)

	eshttp.RegisterDomainRoutes(todo.DOMAIN, mux, store, registry, todo.NewTodoAggregate)

	webServerCfg := web.ServerConfig{
		Router: mux,
	}
	return web.NewHttpServer(webServerCfg)
}
