package memory

import (
	"context"
	"database/sql"
	"sync"

	storedb "github.com/7empestx/GoHTMXToDoList/internal/db/store/sqlc"
)

// MemoryStore implements an in-memory version of the database queries
type MemoryStore struct {
	mu       sync.RWMutex
	tasks    []storedb.Task
	nextID   int32
}

func New() *MemoryStore {
	return &MemoryStore{
		tasks:  make([]storedb.Task, 0),
		nextID: 1,
	}
}

func (m *MemoryStore) GetTasks(ctx context.Context) ([]storedb.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]storedb.Task, len(m.tasks))
	copy(result, m.tasks)
	return result, nil
}

func (m *MemoryStore) AddTask(ctx context.Context, arg storedb.AddTaskParams) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task := storedb.Task{
		ID:          m.nextID,
		Description: arg.Description,
		Completed:   arg.Completed,
		Addedfrom:   arg.Addedfrom,
	}
	m.tasks = append(m.tasks, task)
	m.nextID++
	return nil
}

func (m *MemoryStore) Checked(ctx context.Context, id int32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.tasks {
		if m.tasks[i].ID == id {
			// Toggle the completed status
			m.tasks[i].Completed = sql.NullBool{
				Bool:  !m.tasks[i].Completed.Bool,
				Valid: true,
			}
			break
		}
	}
	return nil
}

func (m *MemoryStore) DeleteTask(ctx context.Context, id int32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.tasks {
		if m.tasks[i].ID == id {
			// Remove task by replacing with last element and truncating
			m.tasks[i] = m.tasks[len(m.tasks)-1]
			m.tasks = m.tasks[:len(m.tasks)-1]
			break
		}
	}
	return nil
}

func (m *MemoryStore) FilterCompletedTasks(ctx context.Context) ([]storedb.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]storedb.Task, 0)
	for _, task := range m.tasks {
		if task.Completed.Valid && task.Completed.Bool {
			result = append(result, task)
		}
	}
	return result, nil
}

func (m *MemoryStore) FilterIncompleteTasks(ctx context.Context) ([]storedb.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]storedb.Task, 0)
	for _, task := range m.tasks {
		if !task.Completed.Valid || !task.Completed.Bool {
			result = append(result, task)
		}
	}
	return result, nil
}
