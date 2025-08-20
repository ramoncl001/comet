package api

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"unicode"

	"github.com/ramoncl001/comet/data"
	"github.com/ramoncl001/comet/ioc"
	"github.com/ramoncl001/comet/middleware"
	"github.com/ramoncl001/comet/rest"
	"github.com/ramoncl001/comet/security"
	"github.com/ramoncl001/comet/security/authentication"
	"github.com/ramoncl001/comet/security/authentication/jwt"
	"gorm.io/gorm"
)

type ApiServer interface {
	MapController(controller interface{})
	UseDatabaseContext(dialector gorm.Dialector, args ...gorm.Option)
	UseMiddleware(m middleware.Middleware)
	AddJWTAuthentication(mg interface{}, provider jwt.JwtProvider, config jwt.JwtConfigurations, userConfig security.UserConfig)
	//UseAuthorization()
	//UseAuthentication()
	Run(addr string) error
}

type apiServer struct {
	ApiServer
	server      *http.ServeMux
	router      *router
	middlewares []middleware.Middleware
}

func CreateServer() ApiServer {
	return &apiServer{
		server: http.NewServeMux(),
		router: newRouter(),
	}
}

func (srv *apiServer) AddJWTAuthentication(mg interface{}, provider jwt.JwtProvider, config jwt.JwtConfigurations, userConfig security.UserConfig) {
	managerType := reflect.TypeOf(mg)
	if managerType.Kind() != reflect.Func {
		panic(managerType.Name() + "is not a SessionManager constructor function")
	}

	ioc.RegisterSingleton(&userConfig)
	ioc.RegisterTransient[security.UserManager](security.NewDefaultUserManager)
	ioc.RegisterSingleton(provider)
	ioc.RegisterSingleton(config)
	ioc.RegisterTransient[authentication.SessionManager](mg)
}

func (srv *apiServer) MapController(controller interface{}) {
	typ := reflect.TypeOf(controller).Out(0)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	name := typ.Name()
	ioc.RegisterKeyedScoped[rest.ControllerBase](controller, name)
	ctrl, err := ioc.ResolveKeyedScoped[rest.ControllerBase](context.Background(), typ.Name())
	if err != nil {
		panic(err)
	}

	srv.router.register(ctrl)
}

func (srv *apiServer) UseMiddleware(m middleware.Middleware) {
	srv.middlewares = append(srv.middlewares, m)
}

func (stc *apiServer) UseDatabaseContext(dialector gorm.Dialector, args ...gorm.Option) {
	ctx := data.NewDatabaseContext(dialector, args...)
	ioc.RegisterSingleton(ctx)
}

func (srv *apiServer) Run(addr string) error {
	fmt.Printf("Running server in %s...\n", addr)
	fmt.Println("Routes:")

	for _, controller := range srv.router.controllers {
		fmt.Printf("[%s]\n", controller.name)
		for _, route := range controller.dynamicRoutes {
			fmt.Printf("[%s]: %s\n", route.Method, route.PathPattern)
		}

		for key := range controller.staticRoutes {
			method := strings.Split(key, ":")[0]
			path := strings.Replace(key, method+":", "", 1)
			fmt.Printf("[%s]: %s\n", method, path+"/")
		}
	}

	middlewares := chain(srv.router.Handle, srv.middlewares...)

	srv.server.Handle("/", middleware.HTTPAdapter(middlewares))
	return http.ListenAndServe(addr, srv.server)
}

func getRouteName(baseName string) string {
	name := strings.ReplaceAll(baseName, "Controller", "")
	var result strings.Builder
	for i, char := range name {
		if unicode.IsUpper(char) {
			if i > 0 && !unicode.IsUpper(rune(name[i-1])) {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}
	name = result.String()
	return "/" + name
}
