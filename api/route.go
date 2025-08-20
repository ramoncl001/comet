package api

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/ramoncl001/comet/ioc"
	"github.com/ramoncl001/comet/rest"
)

type routeHandler struct {
	Function rest.RequestHandler
	Name     string
}

type route struct {
	Method      string
	PathPattern string
	Handler     routeHandler
	PathParts   []string
	ParamNames  []string
}

type controller struct {
	name          string
	basePath      string
	staticRoutes  map[string]*routeHandler
	dynamicRoutes []route
}

type router struct {
	controllers map[string]controller
}

func (r *router) register(ctrl rest.ControllerBase) {
	controllerType := reflect.TypeOf(ctrl)

	controller := controller{
		staticRoutes:  make(map[string]*routeHandler, 0),
		dynamicRoutes: make([]route, 0),
	}

	basePath := ""
	if ctrl.Route() == "" {
		basePath = getRouteName(controllerType.Name())
	} else {
		basePath = ctrl.Route()
	}

	controller.basePath = basePath
	name := strings.ReplaceAll(basePath, "/", "")

	for i := 0; i < controllerType.NumMethod(); i++ {
		method := controllerType.Method(i)
		if !isRequestMethod(method) {
			continue
		}

		invariantName := strings.ToUpper(method.Name)
		methodMap := []rest.RequestMethod{rest.GET, rest.POST, rest.DELETE, rest.PATCH, rest.POST, rest.PUT, rest.LIST}

		for _, prefix := range methodMap {
			if strings.HasPrefix(invariantName, prefix.String()) {
				path := getMethodPath(basePath, method.Name)
				httpMethod := prefix.Method()

				handler := func(req *rest.Request) rest.Response {
					n := controllerType.Name()
					ctrl, err := ioc.ResolveKeyedScoped[rest.ControllerBase](req.Context(), n)
					if err != nil {
						return rest.Error("error getting controller")
					}

					controller := reflect.ValueOf(ctrl)
					request := reflect.ValueOf(req)

					responses := method.Func.Call([]reflect.Value{controller, request})
					result := responses[0].Interface().(rest.Response)
					return result
				}

				if !strings.Contains(path, ":") {
					key := fmt.Sprintf("%s:%s", httpMethod, path)
					controller.staticRoutes[key] = &routeHandler{
						Function: handler,
						Name:     method.Name,
					}
					break
				}

				parts := strings.Split(path, "/")
				params := make([]string, 0)

				for _, part := range parts {
					if strings.HasPrefix(part, ":") {
						params = append(params, part[1:])
					}
				}

				controller.dynamicRoutes = append(controller.dynamicRoutes, route{
					Method:      httpMethod,
					PathPattern: path,
					Handler: struct {
						Function rest.RequestHandler
						Name     string
					}{
						Function: handler,
						Name:     method.Name,
					},
					PathParts:  parts,
					ParamNames: params,
				})

				break
			}
		}
	}

	controller.name = controllerType.Name()
	r.controllers[name] = controller
}

func (r *router) Handle(req *rest.Request) rest.Response {
	path := req.Url.Path

	name := strings.Split(path, "/")[1]
	controller := r.controllers[name]

	str := fmt.Sprintf("%s:%s", req.Method, path)

	var handler *routeHandler
	handler, ok := controller.staticRoutes[str]
	if !ok {
		for _, route := range controller.dynamicRoutes {
			if route.Method != req.Method {
				continue
			}

			if params, ok := r.matchPath(route, path); ok {
				req.PathParams = params
				handler = &route.Handler
			}
		}
	}

	if handler == nil {
		return rest.NotFound()
	}

	ctrl, err := ioc.ResolveKeyedScoped[rest.ControllerBase](req.Context(), controller.name)
	if err != nil {
		return rest.NotFound()
	}

	authorizeMap := make(map[interface{}]rest.AuthorizerFunction)
	config := make([]rest.Policy, 0)

	globalPolicies := ctrl.Policies()["*"]
	if globalPolicies != nil {
		config = append(config, globalPolicies...)
	}

	handlerPolicies := ctrl.Policies()[handler.Name]
	if handlerPolicies != nil {
		config = append(config, handlerPolicies...)
	}

	for _, val := range config {
		authorizeMap[val.Value] = val.Validation
	}

	resultHandler := chainAuthorizations(handler.Function, authorizeMap)
	return resultHandler(req)
}

func (r *router) matchPath(route route, path string) (map[string]string, bool) {
	pathParts := strings.Split(path, "/")

	if len(pathParts) != len(route.PathParts) {
		return nil, false
	}

	params := make(map[string]string)

	for i, part := range route.PathParts {
		if strings.HasPrefix(part, ":") {
			paramName := part[1:]
			params[paramName] = pathParts[i]
		} else if part != pathParts[i] {
			return nil, false
		}
	}

	return params, true
}

func newRouter() *router {
	return &router{
		controllers: make(map[string]controller),
	}
}

func isRequestMethod(method reflect.Method) bool {
	if method.Name[0] < 'A' || method.Name[0] > 'Z' {
		return false
	}

	validPrefixes := []rest.RequestMethod{rest.GET, rest.POST, rest.PUT, rest.PATCH, rest.DELETE, rest.LIST}

	methodName := strings.ToUpper(method.Name)
	matchPrefix := false
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(methodName, prefix.String()) {
			matchPrefix = true
			break
		}
	}

	if !matchPrefix {
		return false
	}

	if method.Type.NumIn() != 2 {
		return false
	}

	reqType := method.Type.In(1)
	if reqType.Kind() != reflect.Ptr {
		return false
	}

	reqType = reqType.Elem()

	if method.Type.NumOut() != 1 {
		return false
	}

	respType := method.Type.Out(0)
	if respType.Kind() == reflect.Ptr {
		respType = respType.Elem()
	}

	if respType != reflect.TypeOf(rest.Response{}) {
		return false
	}

	return reqType.Name() == "Request" && strings.Contains(reqType.PkgPath(), "github.com/ramoncl001/comet")
}

func getMethodPath(basePath, methodName string) string {
	httpMethods := []string{"Get", "Post", "Put", "Delete", "Patch", "List"}
	for _, prefix := range httpMethods {
		if strings.HasPrefix(methodName, prefix) {
			methodName = strings.TrimPrefix(methodName, prefix)
			break
		}
	}

	var result strings.Builder
	for i, char := range methodName {
		if unicode.IsUpper(char) {
			if i > 0 && !unicode.IsUpper(rune(methodName[i-1])) {
				result.WriteRune('/')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}
	route := result.String()

	re := regexp.MustCompile(`(?:^|/)(by|for|of|with)/([^/]+)`)
	route = re.ReplaceAllStringFunc(route, func(match string) string {
		parts := strings.Split(match, "/")
		paramName := parts[len(parts)-1] // El último segmento es el parámetro

		if strings.HasPrefix(match, "/") {
			return "/:" + paramName
		}
		return ":" + paramName
	})

	baseName := strings.ReplaceAll(basePath, "/", "")
	route = strings.ReplaceAll(route, baseName+"/", "")
	route = strings.ReplaceAll(route, baseName, "")

	var completed string
	if len(route) > 0 {
		completed = basePath + "/" + route
	} else {
		completed = basePath
	}

	return completed
}

func chainAuthorizations(handler rest.RequestHandler, authorizers rest.AuthorizationMap) rest.RequestHandler {
	for key, function := range authorizers {
		handler = function(handler, key)
	}
	return handler
}
