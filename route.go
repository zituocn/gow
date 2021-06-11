/*
route.go

like mux  https://github.com/gorilla/mux
@see https://github.com/gorilla/mux/blob/master/route.go

*/

package gow

// routeConfig route config param
type routeConfig struct {

	// if true ignore case on the route
	ignoreCase bool
}

type matcher interface {
	Match(string, *matchValue) bool
}

// matchValue return route struct
type matchValue struct {
	handlers HandlersChain
	params   *Params
	fullPath string
}

type methodTree struct {
	method string
	routes []*Route
}

func (tree methodTree) get(method string) []*Route {
	if tree.method == method {
		return tree.routes
	}
	return nil
}

func (tree methodTree) addRoute(path string, handlers HandlersChain, rc *routeConfig) []*Route {
	route := &Route{
		path:     path,
		handlers: handlers,
	}
	route.addRegexpMatcher(path, rc)
	tree.routes = append(tree.routes, route)
	return tree.routes
}

type Route struct {
	path     string
	handlers HandlersChain
	matchers []matcher
}

// Match implements interface
//	route math
func (r *Route) Match(path string, match *matchValue) bool {
	for _, m := range r.matchers {
		if matched := m.Match(path, match); !matched {
			return false
		}
	}
	match.handlers = r.handlers
	match.fullPath = path
	return true
}

func (r *Route) addMatcher(m matcher) *Route {
	r.matchers = append(r.matchers, m)
	return r
}

func (r *Route) addRegexpMatcher(path string, rc *routeConfig) error {
	rr, err := addRouteRegexp(path, rc)
	if err != nil {
		return err
	}
	r.addMatcher(rr)

	return nil
}
