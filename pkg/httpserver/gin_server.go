package httpserver

import (
	"context"
	"net/http"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/propagation"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 10 * time.Second
)

type Handler interface {
	Method() string
	Path() string
	Handle(c *gin.Context)
}

type Option func(*Server)

func WithHealthCheck() Option {
	return func(s *Server) {
		s.router.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
}

func WithOpenTracing(serviceName string) Option {
	return func(s *Server) {
		s.Use(otelgin.Middleware(
			serviceName,
			otelgin.WithFilter(func(r *http.Request) bool {
				return r.URL.Path != "/health"
			}),
			otelgin.WithPropagators(
				propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{}, propagation.Baggage{},
				),
			),
		))
	}
}

func WithZeroLogger(logger *zerolog.Logger) Option {
	return func(s *Server) {
		s.Use(GinRequestZeroLogger(logger))
	}
}

func WithSwaggerDocs(url string) Option {
	return func(s *Server) {
		s.router.GET(url, ginSwagger.WrapHandler(
			swaggerFiles.Handler,
		))
	}
}

func WithCORS() Option {
	return func(s *Server) {
		s.Use(cors.New(
			cors.Config{
				AllowAllOrigins:  true,
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
				AllowCredentials: true,
			},
		))
	}
}

func WithHandlers(handlers ...Handler) Option {
	return func(s *Server) {
		for _, h := range handlers {
			s.Handle(
				h.Method(),
				path.Join(s.basePath, h.Path()),
				h.Handle,
			)
		}
	}
}

type Server struct {
	basePath string
	port     string

	router *gin.Engine
	srv    *http.Server
}

// Handle adds a new route to the GinServer.
func (s *Server) Handle(method, path string, handler gin.HandlerFunc) {
	s.router.Handle(method, path, handler)
}

// AddMiddleware adds a new middleware to the GinServer.
func (s *Server) Use(middleware gin.HandlerFunc) {
	s.router.Use(middleware)
}

// Run starts the GinServer on the specified port.
func (s *Server) Run(ctx context.Context) error {
	s.srv = &http.Server{
		Addr:              makeAddr(s.port),
		Handler:           s.router,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadTimeout,
		WriteTimeout:      DefaultWriteTimeout,
	}

	go func() {
		<-ctx.Done()
		err := s.srv.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}()

	return s.srv.ListenAndServe()
}

// NewGinServer creates a new GinServer with default healthcheck route and middlewares.
func NewGinServer(basePath, port string, opts ...Option) *Server {
	router := gin.New()
	router.Use(gin.Recovery())

	server := &Server{
		basePath: basePath,
		router:   router,
		port:     port,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func makeAddr(port string) string {
	return ":" + port
}
