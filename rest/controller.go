package rest

type ControllerBase interface {
	Route() string
	Policies() PoliciesConfig
}

type RequestHandler func(*Request) Response
