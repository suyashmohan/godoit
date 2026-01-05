package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	todov1 "github.com/suyashmohan/godoit/gen/todo/v1"
	"github.com/suyashmohan/godoit/gen/todo/v1/todov1connect"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type todoServer struct {
	mu     sync.Mutex
	todos  map[int32]*todov1.Todo
	nextID int32
}

// CreateTodo implements [todov1connect.TodoServiceHandler].
func (t *todoServer) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.CreateTodoResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	todo := &todov1.Todo{
		Id:        t.nextID,
		Text:      req.Text,
		Completed: false,
		CreatedAt: time.Now().String(),
	}
	t.todos[todo.Id] = todo
	t.nextID++

	return &todov1.CreateTodoResponse{
		Todo: todo,
	}, nil
}

// DeleteTodo implements [todov1connect.TodoServiceHandler].
func (t *todoServer) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*todov1.DeleteTodoResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, exists := t.todos[req.Id]
	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("todo not found"))
	}

	delete(t.todos, req.Id)

	return &todov1.DeleteTodoResponse{
		Success: true,
	}, nil
}

// GetTodos implements [todov1connect.TodoServiceHandler].
func (t *todoServer) GetTodos(ctx context.Context, req *todov1.GetTodosRequest) (*todov1.GetTodosResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	todos := make([]*todov1.Todo, 0, len(t.todos))
	for _, todo := range t.todos {
		todos = append(todos, todo)
	}

	return &todov1.GetTodosResponse{
		Todos: todos,
	}, nil
}

// UpdateTodo implements [todov1connect.TodoServiceHandler].
func (t *todoServer) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.UpdateTodoResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	todo, exists := t.todos[req.Id]
	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("todo not found"))
	}

	todo.Completed = req.Completed

	return &todov1.UpdateTodoResponse{
		Todo: todo,
	}, nil
}

func newTodoServer() *todoServer {
	return &todoServer{
		todos:  make(map[int32]*todov1.Todo),
		nextID: 1,
	}
}

func main() {
	mux := http.NewServeMux()

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	server := newTodoServer()
	path, handler := todov1connect.NewTodoServiceHandler(server)
	mux.Handle(path, corsMiddleware(handler))

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{}))
}
