package comet

import (
	"context"

	"github.com/ramoncl001/comet/api"
	"github.com/ramoncl001/comet/ioc"
	"github.com/ramoncl001/comet/log"
	"github.com/ramoncl001/comet/middleware"
	"github.com/ramoncl001/comet/rest"
)

// ApiServer represents the main API server instance.
// Handles HTTP requests, routing, middleware, and application lifecycle.
type ApiServer api.ApiServer

// ControllerBase provides the foundation for all API controllers.
// Offers helper methods for request handling and response generation.
type ControllerBase rest.ControllerBase

// Request encapsulates the incoming HTTP request with convenient methods
// to access parameters, headers, body content, and other request data.
type Request rest.Request

// Response represents the HTTP response to be sent to the client.
// Provides methods to set status codes, headers, and response body content.
type Response rest.Response

// RequestHandler is a function type that processes HTTP requests and generates responses.
// The fundamental building block for defining API endpoints and handlers.
type RequestHandler = func(r *rest.Request) rest.Response

// Middleware is a function that intercepts and processes HTTP requests
// before they reach the main handler, enabling cross-cutting concerns.
type Middleware = func(next rest.RequestHandler) rest.RequestHandler

// PoliciesConfig is a map with all controller policies
// configuration, such as Role, Permission or custom policies
type PoliciesConfig rest.PoliciesConfig

// NewServer creates and returns a new instance of the API server.
// This is the entry point for initializing the Comet framework application.
func NewServer() ApiServer {
	return api.CreateServer()
}

// RegisterTransient registers a transient dependency in the IoC container.
// A new instance is created every time the dependency is resolved.
func RegisterTransient[T any](provider interface{}) {
	ioc.RegisterTransient[T](provider)
}

// RegisterKeyedTransient registers a keyed transient dependency in the IoC container.
// Allows multiple implementations of the same interface to be registered with different keys.
func RegisterKeyedTransient[T any](provider, key interface{}) {
	ioc.RegisterKeyedTransient[T](provider, key)
}

// RegisterScoped registers a scoped dependency in the IoC container.
// The same instance is reused within the same context/request scope.
func RegisterScoped[T any](provider interface{}) {
	ioc.RegisterScoped[T](provider)
}

// RegisterKeyedScoped registers a keyed scoped dependency in the IoC container.
// Provides scoped resolution with key-based implementation selection.
func RegisterKeyedScoped[T any](provider, key interface{}) {
	ioc.RegisterKeyedScoped[T](provider, key)
}

// RegisterSingleton registers a singleton dependency in the IoC container.
// A single instance is created and reused for the entire application lifetime.
func RegisterSingleton[T any](instance T) {
	ioc.RegisterSingleton(instance)
}

// RegisterKeyedSingleton registers a keyed singleton dependency in the IoC container.
// Enables singleton resolution with key-based implementation selection.
func RegisterKeyedSingleton[T any](instance T, key interface{}) {
	ioc.RegisterKeyedSingleton(instance, key)
}

// ResolveTransient resolves a transient dependency from the IoC container.
// Returns a new instance of the requested type with all dependencies injected.
func ResolveTransient[T any](ctx context.Context) (T, error) {
	return ioc.ResolveTransient[T](ctx)
}

// ResolveKeyedTransient resolves a keyed transient dependency from the IoC container.
// Returns a new instance of the requested type based on the provided key.
func ResolveKeyedTransient[T any](ctx context.Context, key interface{}) (T, error) {
	return ioc.ResolveKeyedTransient[T](ctx, key)
}

// ResolveScoped resolves a scoped dependency from the IoC container.
// Returns the same instance within the same context/request scope.
func ResolveScoped[T any](ctx context.Context) (T, error) {
	return ioc.ResolveScoped[T](ctx)
}

// ResolveKeyedScoped resolves a keyed scoped dependency from the IoC container.
// Returns a scoped instance based on the provided key identifier.
func ResolveKeyedScoped[T any](ctx context.Context, key interface{}) (T, error) {
	return ioc.ResolveKeyedScoped[T](ctx, key)
}

// ResolveSingleton resolves a singleton dependency from the IoC container.
// Returns the single shared instance of the requested type.
func ResolveSingleton[T any](ctx context.Context) (T, error) {
	return ioc.ResolveSingleton[T](ctx)
}

// ResolveKeyedSingleton resolves a keyed singleton dependency from the IoC container.
// Returns the singleton instance associated with the provided key.
func ResolveKeyedSingleton[T any](ctx context.Context, key interface{}) (T, error) {
	return ioc.ResolveKeyedSingleton[T](ctx, key)
}

// LoggerFromContext retrieves a logger instance from the context.
// Provides consistent logging throughout the application with context-aware capabilities.
func LoggerFromContext(ctx context.Context) log.Logger {
	return log.FromContext(ctx)
}

// RequestLogging middleware automatically logs incoming HTTP requests
// and responses with relevant timing and metadata information.
var RequestLogging = middleware.RequestLogging

// Recover middleware captures and handles panics gracefully,
// preventing server crashes and providing structured error responses.
var Recover = middleware.Recover

// RequestID middleware automatically generates and assigns unique identifiers
// to each incoming request for improved tracing and debugging capabilities.
var RequestID = middleware.RequestID
