package router

import (
	"fmt"
	"net/http"
	"runtime"

	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func NewRouter() (*http.Server, *mux.Router) {
	rtr := mux.NewRouter()
	rtr.StrictSlash(false)
	var allowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	var listenURL = ":8080"
	addAppRoutes(rtr.PathPrefix("/v1").Subrouter())
	n := negroni.New(NewPanicHandler(), negroni.NewLogger(), negroni.NewStatic(http.Dir("public")))
	n.UseHandler(rtr)
	corsHeaderOptions := handlers.AllowedHeaders([]string{"Access-Control-Max-Age",
		"Origin", "X-Requested-With", "Content-Type", "Accept", "Content-Length",
		"Accept-Encoding", "X-CSRF-Token", "Authorization"})
	corsOriginOptions := handlers.AllowedOrigins(allowedOrigins)
	corsMethodOptions := handlers.AllowedMethods([]string{"GET", "HEAD",
		"POST", "PUT", "DELETE", "OPTIONS"})
	corsCredentialsOptions := handlers.AllowCredentials()
	srv := &http.Server{Addr: listenURL}
	srv.Handler = handlers.CORS(corsHeaderOptions, corsOriginOptions, corsMethodOptions, corsCredentialsOptions)(n)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()
	return srv, rtr
}

type PanicHandler struct {
	PrintStack       bool
	ErrorHandlerFunc func(interface{})
	StackAll         bool
	StackSize        int
}

func NewPanicHandler() *PanicHandler {
	return &PanicHandler{
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (fn *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, fn.StackSize)
			stack = stack[:runtime.Stack(stack, fn.StackAll)]

			if fn.PrintStack {
				handleError(w, fmt.Errorf("%v", err), "server error", 500)
				handleError(w, fmt.Errorf("%v", string(stack)), "server error", 500)
			}

			if fn.ErrorHandlerFunc != nil {
				func() {
					defer func() {
						if err := recover(); err != nil {
							handleError(w, fmt.Errorf("%v", err), "server error", 500)
							handleError(w, fmt.Errorf("%v", string(stack)), "server error", 500)
						}
					}()
					fn.ErrorHandlerFunc(err)
				}()
			}
		}
	}()
	next(w, r)
}

func handleError(w http.ResponseWriter, inputErr error, message string, statusCode int) {
	fmt.Println(inputErr)
	http.Error(w, message, statusCode)
}
