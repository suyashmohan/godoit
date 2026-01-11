package services

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/suyashmohan/godoit/gen/database"
	todov1 "github.com/suyashmohan/godoit/gen/todo/v1"
)

type todoService struct {
	queries *database.Queries
}

func NewTodoService(queries *database.Queries) *todoService {
	return &todoService{
		queries: queries,
	}
}

// CreateTodo implements [todov1connect.TodoServiceHandler].
func (t *todoService) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.CreateTodoResponse, error) {
	todo, err := t.queries.CreateTodo(ctx, database.CreateTodoParams{
		Text:      req.Text,
		Completed: false,
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &todov1.CreateTodoResponse{
		Todo: &todov1.Todo{
			Id:        todo.ID,
			Text:      todo.Text,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt.Time.String(),
		},
	}, nil
}

// DeleteTodo implements [todov1connect.TodoServiceHandler].
func (t *todoService) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*todov1.DeleteTodoResponse, error) {
	err := t.queries.DeleteTodo(ctx, req.Id)
	if err != nil {
		log.Fatalln(err)
	}

	return &todov1.DeleteTodoResponse{
		Success: true,
	}, nil
}

// GetTodos implements [todov1connect.TodoServiceHandler].
func (t *todoService) GetTodos(ctx context.Context, req *todov1.GetTodosRequest) (*todov1.GetTodosResponse, error) {
	todos, err := t.queries.ListTodos(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	todosRes := make([]*todov1.Todo, len(todos))
	for idx, todo := range todos {
		todosRes[idx] = &todov1.Todo{
			Id:        todo.ID,
			Text:      todo.Text,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt.Time.String(),
		}
	}

	return &todov1.GetTodosResponse{
		Todos: todosRes,
	}, nil
}

// UpdateTodo implements [todov1connect.TodoServiceHandler].
func (t *todoService) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.UpdateTodoResponse, error) {
	todo, err := t.queries.UpdateTodo(ctx, database.UpdateTodoParams{
		ID:        req.Id,
		Completed: req.Completed,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &todov1.UpdateTodoResponse{
		Todo: &todov1.Todo{
			Id:        todo.ID,
			Text:      todo.Text,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt.Time.String(),
		},
	}, nil
}
