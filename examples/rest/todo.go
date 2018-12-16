package main

type todo struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

func newTodoStore() *todoStore {
	return &todoStore{
		todos: []*todo{
			&todo{ID: "t1", Title: "Buy milk"},
			&todo{ID: "t2", Title: "Feed gother"},
			&todo{ID: "t3", Title: "Clean room"},
			&todo{ID: "t4", Title: "Be happy"},
		},
	}
}

type todoStore struct {
	todos []*todo
}

func (s *todoStore) GetAll() ([]*todo, error) {
	return s.todos, nil
}

func (s *todoStore) GetByID(todoID string) (*todo, error) {
	for _, t := range s.todos {
		if t.ID == todoID {
			return t, nil
		}
	}
	return nil, nil
}

func (s *todoStore) AddTodo(newTodo *todo) error {
	s.todos = append(s.todos, newTodo)
	return nil
}

func (s *todoStore) DeleteTodo(todoID string) error {
	for i, t := range s.todos {
		if t.ID == todoID {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			break
		}
	}
	return nil
}
