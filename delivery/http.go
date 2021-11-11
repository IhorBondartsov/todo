package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"
	"to-do/api"
	"to-do/app"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	ToDoIDParam = "todoid"
)

type HTTPConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration

	InitProfiling bool
}

func (h *HTTPConfig) Validate() error {
	errs := []error{}
	if h.Host == "" {
		errs = append(errs, errors.New("empty http host"))
	}

	if h.Port < 1 || h.Port > 65535 {
		errs = append(errs, errors.New("wrong http port"))
	}

	if len(errs) > 0 {
		return errors.Errorf("http cfg errors: %v", errs)
	}

	return nil
}

type httpService struct {
	HTTPConfig
	todoService *app.ToDoService
	router      *httprouter.Router
}

func NewHTTPService(cfg HTTPConfig, todoService *app.ToDoService) *httpService {
	service := httpService{
		HTTPConfig:  cfg,
		todoService: todoService,
		router:      httprouter.New(),
	}
	if cfg.InitProfiling {
		service.pprofHandlers("/debug/pprof")
	}
	service.registerRoutes()
	return &service
}

func (s *httpService) Run() {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), s.router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (s *httpService) registerRoutes() {
	s.router.GET("/todo/:todoid", logMiddleware(s.getToDo))
	s.router.PUT("/todo", logMiddleware(s.updateToDo))
	s.router.POST("/todo", logMiddleware(s.createToDo))
	s.router.DELETE("/todo/:todoid", logMiddleware(s.deleteToDo))
}

func (s *httpService) pprofHandlers(path string) {
	s.router.GET(path+"/cmdline", wrapHandlerFunc(pprof.Cmdline))
	s.router.GET(path+"/profile", wrapHandlerFunc(pprof.Profile))
	s.router.GET(path+"/symbol", wrapHandlerFunc(pprof.Symbol))
	s.router.GET(path+"/goroutine", wrapHandler(pprof.Handler("goroutine")))
	s.router.GET(path+"/heap", wrapHandler(pprof.Handler("heap")))
	s.router.GET(path+"/threadcreate", wrapHandler(pprof.Handler("threadcreate")))
	s.router.GET(path+"/block", wrapHandler(pprof.Handler("block")))
	s.router.GET(path+"/trace", wrapHandler(pprof.Handler("trace")))
}

func wrapHandlerFunc(handler http.HandlerFunc) httprouter.Handle {
	handle := func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ctx := req.Context()
		ctx = context.WithValue(ctx, "params", p)
		req = req.WithContext(ctx)
		handler.ServeHTTP(w, req)
	}
	return handle
}

func wrapHandler(handler http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ctx := req.Context()
		ctx = context.WithValue(ctx, "params", p)
		req = req.WithContext(ctx)
		handler.ServeHTTP(w, req)
	}
}

func (s *httpService) getToDo(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := req.Context()
	var todoIDStr string
	for _, v := range params {
		if v.Key == ToDoIDParam {
			todoIDStr = v.Value
			break
		}
	}

	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo, err := s.todoService.GetTodo(ctx, todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if todo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (s *httpService) updateToDo(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := req.Context()
	var newTodo api.ToDo
	if err := json.NewDecoder(req.Body).Decode(&newTodo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(newTodo.Message) == 0 {
		http.Error(w, "message cant be empty", http.StatusNoContent)
		return
	}

	err := s.todoService.UpdateToDo(ctx, newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *httpService) createToDo(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := req.Context()
	var newTodo api.ToDo
	if err := json.NewDecoder(req.Body).Decode(&newTodo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := s.todoService.CreateToDo(ctx, newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *httpService) deleteToDo(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := req.Context()
	var todoIdStr string
	for _, v := range params {
		if v.Key == ToDoIDParam {
			todoIdStr = v.Value
			break
		}
	}

	todoId, err := strconv.ParseInt(todoIdStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.todoService.DeleteTodo(ctx, todoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
