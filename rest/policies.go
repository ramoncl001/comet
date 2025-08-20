package rest

type AuthorizerFunction = func(RequestHandler, interface{}) RequestHandler

type AuthorizationMap map[interface{}]AuthorizerFunction

type PoliciesConfig map[string][]Policy

type Policy struct {
	Validation AuthorizerFunction
	Value      interface{}
}

func Authorize(fn AuthorizerFunction, val interface{}) Policy {
	return Policy{
		Validation: fn,
		Value:      val,
	}
}
