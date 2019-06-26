package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/alexzimmer96/gqlgen-example/graphql"
	"github.com/alexzimmer96/gqlgen-example/repository"
	"github.com/alexzimmer96/gqlgen-example/service"
	"github.com/dimiro1/health"
	"github.com/gorilla/websocket"
	"github.com/muesli/cache2go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	// Initializing Database, Repositories and Services here
	db := cache2go.Cache("example")
	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)

	// Create a new Router and Register the GraphQL-Resolver
	res := graphql.NewResolver(articleService)

	// Some GraphQL-Configuration
	graphqlConfig := graphql.NewExecutableSchema(graphql.Config{Resolvers: res})
	websocketUpgrader := handler.WebsocketUpgrader(websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}},
	)
	// Set ping-time for sockets to 5 seconds to prevent disconnects
	websocketKeepalive := handler.WebsocketKeepAliveDuration(time.Second * 5)

	graphqlHandler := handler.GraphQL(graphqlConfig, websocketUpgrader, websocketKeepalive)
	http.Handle("/query", graphqlHandler)

	// Adding Playground, maybe adding a debug-mode switch later
	playgroundHandler := handler.Playground("GraphQL", "/query")
	http.Handle("/playground", playgroundHandler)

	// Adding "/metrics" endpoint for prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Adding "/status" endpoint for health-checking
	healthHandler := health.NewHandler()
	http.Handle("/status", healthHandler)
	http.HandleFunc("/_stack", getStackTraceHandler)

	// Finally starting the HTTP-Server
	startHttpServer(1337)
}

// Starting a HTTP-Server using the router object an a given port
// Handles graceful-shutdowns
func startHttpServer(port int) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()
	logrus.Info(fmt.Sprintf("server started and is available on %d", port))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logrus.Info("received shutdown signal. Trying to shutdown gracefully")

	// Save a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		logrus.WithError(err).Error("failure while shutting down gracefully")
	} else {
		logrus.Info("shutdown completed")
	}
	os.Exit(0)
}

func getStackTraceHandler(w http.ResponseWriter, r *http.Request) {
	stack := debug.Stack()
	w.Write(stack)
	pprof.Lookup("goroutine").WriteTo(w, 2)
}
