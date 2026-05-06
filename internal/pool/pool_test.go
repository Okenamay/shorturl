package pool

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockObject - простая структура для тестов, реализующая интерфейс Resetter
type mockObject struct {
	id    int
	value string
}

// Reset очищает состояние объекта
func (m *mockObject) Reset() {
	m.value = ""
}

func TestPool_GetPut(t *testing.T) {
	counter := 0
	// Фабрика создает объекты с инкрементальным ID
	factory := func() *mockObject {
		counter++
		return &mockObject{id: counter, value: "initial"}
	}

	// Создаем пул с вместимостью 1
	p := New(factory, 1)

	// 1. Получаем объект (пул пуст -> создается новый #1)
	obj1 := p.Get()
	assert.Equal(t, 1, obj1.id)
	assert.Equal(t, "initial", obj1.value)

	// Модифицируем объект, чтобы проверить работу Reset
	obj1.value = "dirty"

	// 2. Возвращаем объект в пул (он должен сброситься)
	p.Put(obj1)
	assert.Equal(t, "", obj1.value, "Value should be reset")

	// 3. Получаем объект снова (должен вернуться тот же #1, так как он был в пуле)
	obj2 := p.Get()
	assert.Equal(t, 1, obj2.id)
	// Проверяем, что это тот же указатель
	assert.Same(t, obj1, obj2)

	// 4. Получаем еще один объект (пул пуст -> создается новый #2)
	obj3 := p.Get()
	assert.Equal(t, 2, obj3.id)

	// 5. Возвращаем оба. Пул размером 1, поэтому один из них отбросится.
	p.Put(obj2)
	p.Put(obj3)

	// Проверяем, что в канале только 1 элемент (максимальная вместимость)
	assert.Equal(t, 1, len(p.store))
}

func TestPool_Concurrency(t *testing.T) {
	factory := func() *mockObject {
		return &mockObject{}
	}
	// Пул достаточного размера
	p := New(factory, 50)

	var wg sync.WaitGroup
	iterations := 1000

	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			// Активно берем и возвращаем объекты в параллельных горутинах
			obj := p.Get()
			obj.value = "working"
			p.Put(obj)
		}()
	}

	wg.Wait()
	// Если тест не упал с deadlock или panic, значит пул потокобезопасен
}
