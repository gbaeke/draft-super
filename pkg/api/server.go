package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	dapr "github.com/dapr/go-sdk/client"
	_ "github.com/gbaeke/super-api/pkg/api/docs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
)

// @title Super API
// @version 0.1
// @description Super API

// @contact.name Source Code
// @contact.url https://github.com/gbaeke/super-api

// @host localhost:8080
// @BasePath /
// @schemes http https

//Config API configuration via viper
type Config struct {
	Welcome    string
	Port       int
	Log        bool
	Timeout    time.Duration
	Daprport   int
	Statestore string
	Pubsub     string
}

//Server struct
type Server struct {
	config     *Config
	logger     *zap.SugaredLogger
	router     *mux.Router
	daprClient dapr.Client
}

//NewServer creates new server
func NewServer(config *Config, logger *zap.SugaredLogger) (*Server, error) {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	srv := &Server{
		config:     config,
		logger:     logger,
		router:     mux.NewRouter(),
		daprClient: client,
	}

	return srv, nil
}

//SetupRoutes sets up routes
func (s *Server) setupRoutes() {
	s.logger.Infow("Enabling /metrics route")
	s.router.Handle("/metrics", promhttp.Handler())

	s.logger.Infow("Enabling /healthz and /readyz routes")
	s.router.HandleFunc("/healthz", s.healthz)
	s.router.HandleFunc("/readyz", s.readyz)

	s.logger.Infow("Enabling /source route")
	s.router.HandleFunc("/source", s.sourceIpHandler)

	s.logger.Infow("Enabling /call route")
	s.router.HandleFunc("/call", s.callMethod)

	s.logger.Infow("Enabling /flaky route")
	s.router.HandleFunc("/flaky", s.flakyHandler)

	s.logger.Infow("Enabling Dapr routes")
	s.router.HandleFunc("/savestate", s.saveState)
	s.router.HandleFunc("/readstate", s.readState)
	s.router.HandleFunc("/dapr/subscribe", s.daprSubScribe)
	s.router.HandleFunc("/myroute", s.myRoute)
	s.router.HandleFunc("/mqtt", s.mqtt)

	s.logger.Infow("Enabling index route")
	s.router.HandleFunc("/", s.indexHandler)

	s.logger.Infow("Enabling auth route")
	s.router.HandleFunc("/auth", s.authHandler)

	s.logger.Infow("Enabling /swagger.json route")
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	s.router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if err != nil {
			s.logger.Error("swagger error", zap.Error(err), zap.String("path", "/swagger.json"))
		}
		w.Write([]byte(doc))
	})
}

func (s *Server) setupMiddlewares() {
	if s.config.Log {
		// only log requests when --log is set
		s.router.Use(s.loggingMiddleware)
	}
}

//StartServer starts http server
func (s *Server) StartServer() {

	s.setupRoutes()
	s.setupMiddlewares()

	srv := &http.Server{
		Addr:    ":" + fmt.Sprint(s.config.Port),
		Handler: s.router,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.logger.Infow("starting web server",
		zap.Int("port", s.config.Port),
	)

	//graceful shutdown - run server in goroutine and handle SIGINT & SIGTERM
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			s.logger.Infow("server stopped",
				zap.Error(err),
			)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// block wait for signal
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()
	srv.Shutdown(ctx)

	s.logger.Infow("server shutting down")
	os.Exit(0)
}
