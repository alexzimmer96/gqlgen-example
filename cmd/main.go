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

	// Creating handler for application-logic and add the endpoints
	applicationHandler := http.NewServeMux()
	graphqlHandler := getGraphQLHandler(articleService)
	applicationHandler.Handle("/query", graphqlHandler)
	playgroundHandler := handler.Playground("GraphQL", "/query")
	applicationHandler.Handle("/playground", playgroundHandler)

	// Creating handler for monitoring and add the endpoints
	monitoringHandler := http.NewServeMux()
	monitoringHandler.Handle("/metrics", promhttp.Handler())
	monitoringHandler.Handle("/status", health.NewHandler())

	// Finally starting the HTTP-Server
	startHttpServer(1337, 1338, applicationHandler, monitoringHandler)
}

// Starting a HTTP-Server using the router object an a given port
// Handles graceful-shutdowns
func startHttpServer(appPort, monitoringPort int, appHandler, monitoringHandler http.Handler) {
	appSrv := startServerInstance(appPort, appHandler)
	monitoringSrv := startServerInstance(monitoringPort, monitoringHandler)
	logrus.Info(fmt.Sprintf("started application-server on port %d. Monitoring-server is available on port %d", appPort, monitoringPort))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logrus.Info("received shutdown signal. Trying to shutdown gracefully")

	// Save a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err, err2 := appSrv.Shutdown(ctx), monitoringSrv.Shutdown(ctx)
	if err != nil {
		logrus.WithError(err).Error("failure while shutting down application-server gracefully")
	} else if err2 != nil {
		logrus.WithError(err).Error("failure while shutting down monitoring-server gracefully")
	} else {
		logrus.Info("shutdown completed")
	}
	os.Exit(0)
}

func startServerInstance(port int, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Error(err)
		}
	}()
	return srv
}

func getGraphQLHandler(articleService service.IArticleService) http.HandlerFunc {
	res := graphql.NewResolver(articleService)

	graphqlConfig := graphql.NewExecutableSchema(graphql.Config{Resolvers: res})
	websocketUpgrader := handler.WebsocketUpgrader(websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}},
	)

	websocketKeepalive := handler.WebsocketKeepAliveDuration(time.Second * 5)
	return handler.GraphQL(graphqlConfig, websocketUpgrader, websocketKeepalive)
}
