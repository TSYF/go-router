package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Handler func(http.ResponseWriter, *RequestData)
type PreHandler func(http.ResponseWriter, *RequestData)

type HttpHandler struct{}

type HttpMethod string

type RequestData struct {
	Params map[string]string
	http.Request
}

type Router struct {
	middleware map[string]map[string]*PreHandler `validate:"required"`
	handlers map[string]map[string]*PreHandler `validate:"required"`
	parametizedRoutes map[string]map[string]string
	paramRegex *regexp.Regexp `validate:"required"`
	/*
		? handlers map[string]map[string]PreHandler
		? Example: { "/products": { "POST": PreHandler } }
	*/
}

func NewRouter() *Router {
	return &Router{
		middleware: make(map[string]map[string]*PreHandler),
		handlers: make(map[string]map[string]*PreHandler),
		parametizedRoutes: make(map[string]map[string]string),
		paramRegex: regexp.MustCompile("\\{\\w+\\}"),
	}
}

func (r *Router) Get(prefix string, handler Handler) {
	GET := http.MethodGet
	if _, ok := r.handlers[GET]; !ok {
        r.handlers[GET] = make(map[string]*PreHandler)
    }
	var preHandler PreHandler = r.makePreHandler(GET, prefix, handler)
	r.handlers[GET][prefix] = &preHandler
	if r.hasParams(prefix) {
		if _, ok := r.parametizedRoutes[GET]; !ok {
			r.parametizedRoutes[GET] = make(map[string]string)
		}
		r.parametizedRoutes[GET][r.getPrefixBeforeParams(prefix)] = prefix
	}
}

func (r *Router) Post(prefix string, handler Handler) {
	POST := http.MethodPost
	if _, ok := r.handlers[POST]; !ok {
        r.handlers[POST] = make(map[string]*PreHandler)
    }
	var preHandler PreHandler = r.makePreHandler(POST, prefix, handler)
	r.handlers[POST][prefix] = &preHandler
	if r.hasParams(prefix) {
		if _, ok := r.parametizedRoutes[POST]; !ok {
			r.parametizedRoutes[POST] = make(map[string]string)
		}
		r.parametizedRoutes[POST][r.getPrefixBeforeParams(prefix)] = prefix
	}
}

func (r *Router) Put(prefix string, handler Handler) {
	PUT := http.MethodPut
	if _, ok := r.handlers[PUT]; !ok {
        r.handlers[PUT] = make(map[string]*PreHandler)
    }
	var preHandler PreHandler = r.makePreHandler(PUT, prefix, handler)
	r.handlers[PUT][prefix] = &preHandler
	if r.hasParams(prefix) {
		if _, ok := r.parametizedRoutes[PUT]; !ok {
			r.parametizedRoutes[PUT] = make(map[string]string)
		}
		r.parametizedRoutes[PUT][r.getPrefixBeforeParams(prefix)] = prefix
	}
}

func (r *Router) Patch(prefix string, handler Handler) {
	PATCH := http.MethodPatch
	if _, ok := r.handlers[PATCH]; !ok {
        r.handlers[PATCH] = make(map[string]*PreHandler)
    }
	var preHandler PreHandler = r.makePreHandler(PATCH, prefix, handler)
	r.handlers[PATCH][prefix] = &preHandler
	if r.hasParams(prefix) {
		if _, ok := r.parametizedRoutes[PATCH]; !ok {
			r.parametizedRoutes[PATCH] = make(map[string]string)
		}
		r.parametizedRoutes[PATCH][r.getPrefixBeforeParams(prefix)] = prefix
	}
}

func (r *Router) Delete(prefix string, handler Handler) {
	DELETE := http.MethodDelete
	if _, ok := r.handlers[DELETE]; !ok {
        r.handlers[DELETE] = make(map[string]*PreHandler)
    }
	var preHandler PreHandler = r.makePreHandler(DELETE, prefix, handler)
	r.handlers[DELETE][prefix] = &preHandler
	if r.hasParams(prefix) {
		if _, ok := r.parametizedRoutes[DELETE]; !ok {
			r.parametizedRoutes[DELETE] = make(map[string]string)
		}
		r.parametizedRoutes[DELETE][r.getPrefixBeforeParams(prefix)] = prefix
	}
}

func (r Router) makePreHandler(method string, prefix string, handler Handler) PreHandler {
	
	return func(res http.ResponseWriter, req *RequestData) {
		if req.Method != method {
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed) 
			log.Fatalln("Unallowed method!")
		}

		if r.hasParams(prefix) {
			req.Params = r.getParams(prefix, req.URL.Path)
		}
		
		fmt.Println("handling")
		handler(res, req)
	}
}

func (r Router) getParams(prefix string, path string) map[string]string {
	params := make(map[string]string)
	re := r.paramRegex

	for _, match := range re.FindAllString(prefix, -1) {
		matchIndex := strings.Index(prefix, match)
		matchToNextSlash := strings.Index(path[matchIndex:], "/")
		var nextSlashIndex int
		if matchToNextSlash == -1 {
			nextSlashIndex = len(path)
		} else {
			nextSlashIndex = matchIndex + matchToNextSlash 
		}
		trimmedMatch := match[1:len(match)-1]
		value := path[matchIndex:nextSlashIndex]
		params[trimmedMatch] = value
		fmt.Printf("============\n")
		fmt.Printf("params: %s\n", params)
		fmt.Printf("============\n")
	}
	
	return params
}

func (r Router) getPrefixBeforeParams(prefix string) string {
	re := r.paramRegex
	matchIndex := re.FindStringIndex(prefix)

	prefixBeforeParams := prefix[:matchIndex[0] - 1]
	return prefixBeforeParams
}

// This function will return the key of the route with params
// if it exists, otherwise it will return an error
func (r Router) getParametizedRouteKey(method string, path string) (string, error) {
	routeElements := strings.Split(path, "/")

	if len(routeElements) <= 1 {
		return "", errors.New("Route not found")
	}
	
	firstElements := routeElements[:len(routeElements) - 1]
	p := strings.Join(firstElements, "/")

	if key, ok := r.parametizedRoutes[method][path]; ok {
		return key, nil
	}

	return r.getParametizedRouteKey(method, p)
}


func (r Router) hasParams(route string) bool {
	return r.paramRegex.MatchString(route)
}

func (r Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	prefix := req.RequestURI
	method := req.Method
	
	request := &RequestData{
		Params: make(map[string]string),
		Request: *req,
	}

	for key, middleware := range r.middleware[method] {
		if strings.HasPrefix(prefix, key) {
			(*middleware)(res, request)
		}
	}

	if handler, ok := r.handlers[method][prefix]; ok {
		(*handler)(res, request)
	} else {
		fmt.Println("============")
		fmt.Println(r.parametizedRoutes[method])
		fmt.Println("============")
		if key, err := r.getParametizedRouteKey(method, prefix); err == nil {
			fmt.Println("parameterized route key")
			fmt.Println(key)
			if handler, ok := r.handlers[method][key]; ok {
				fmt.Println("parameterized route handler")
				(*handler)(res, request)
			}
		}
	}

	defer request.Body.Close()
}

func (r Router) Listen(prefix string, port int, host string, confirmationPrintable string) { 
	defaultString(&prefix, "/")
	defaultInt(&port, 8000)
	defaultString(&host, "127.0.0.1")

	fmt.Println(confirmationPrintable)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func (r *Router) Use(method string, path string, middleware Handler) {
	if _, ok := r.handlers[method]; !ok {
        r.middleware[method] = make(map[string]*PreHandler)
    }
	preHandler := r.makePreHandler(method, path, middleware)
	r.middleware[method][path] = &preHandler
}

func defaultInt(variable *int, defaultValue int) { 
	if *variable == 0 {
		*variable = defaultValue
	}
}

func defaultString(variable *string, defaultValue string) {
	if *variable == "" {
		*variable = defaultValue
	}
}
