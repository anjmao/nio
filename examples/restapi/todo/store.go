package todo

type todo struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type TodoStore interface {
	GetAll() ([]*todo, error)
	GetByID(todoID string) (*todo, error)
	Add(newTodo *todo) error
	Delete(todoID string) error
}

func NewTodoStore() TodoStore {
	return &todoStore{
		todos: []*todo{
			{ID: "t1", Title: "Buy milk"},
			{ID: "t2", Title: "Feed gother"},
			{ID: "t3", Title: "Clean room"},
			{ID: "t4", Title: "Be happy"},
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

func (s *todoStore) Add(newTodo *todo) error {
	s.todos = append(s.todos, newTodo)
	return nil
}

func (s *todoStore) Delete(todoID string) error {
	for i, t := range s.todos {
		if t.ID == todoID {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			break
		}
	}
	return nil
}
