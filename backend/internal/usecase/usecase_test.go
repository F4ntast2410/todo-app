package usecase

// import (
// 	"context"
// 	"testing"
// )

// type MockTaskRepository struct {
// 	// Сюда мы можем передать кастомную ошибку для тестов на сбой БД
// 	ErrToReturn error
// 	// Фейковый ID, который "база" присвоит задаче
// 	FakeID int
// }

// // Теперь метод идеально соответствует твоему интерфейсу!
// func (m *MockTaskRepository) Save(ctx context.Context, task *Task) error {
// 	if m.ErrToReturn != nil {
// 		return m.ErrToReturn
// 	}

// 	// Имитируем RETURNING id: перезаписываем ID в переданном объекте
// 	if m.FakeID != 0 {
// 		task.ID = m.FakeID
// 	} else {
// 		task.ID = 1 // дефолтный ID для тестов
// 	}

// 	return nil
// }

// // Заглушки для остальных методов (чтобы интерфейс удовлетворялся)
// func (m *MockTaskRepository) GetAll(ctx context.Context) ([]Task, error)                { return nil, nil }
// func (m *MockTaskRepository) UpdateStatus(ctx context.Context, id int, done bool) error { return nil }
// func (m *MockTaskRepository) Delete(ctx context.Context, id int) error                  { return nil }
// func (m *MockTaskRepository) GetByUserID(ctx context.Context, userID int64) ([]Task, error) {
// 	return nil, nil
// }

// func TestCreateTask_Validation(t *testing.T) {
// 	// 1. Инициализируем окружение
// 	mockRepo := &MockTaskRepository{FakeID: 100}
// 	uc := &TaskUsecaseImpl{Repo: mockRepo}

// 	// 2. Описываем тестовые случаи (таблицу)
// 	tests := []struct {
// 		name    string // Имя теста для вывода в консоль
// 		title   string // Входные данные
// 		wantErr bool   // Ждем ли мы ошибку?
// 	}{
// 		{
// 			name:    "Успешное создание",
// 			title:   "Купить хлеб",
// 			wantErr: false,
// 		},
// 		{
// 			name:    "Пустая строка должна выдавать ошибку",
// 			title:   "",
// 			wantErr: true,
// 		},
// 		// Сюда можно легко докидывать новые кейсы за одну секунду!
// 	}

// 	// 3. Бежим по таблице циклом
// 	for _, tt := range tests {
// 		// t.Run создает изолированный подтест (subtest)
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := uc.CreateTask(context.Background(), tt.title)

// 			// Проверяем, совпадает ли наличие ошибки с нашими ожиданиями
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
