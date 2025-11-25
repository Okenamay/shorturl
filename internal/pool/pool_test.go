package pool

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStruct - это тестовая структура, реализующая Resetter. Важно, что
// *TestStruct реализует интерфейс, а не TestStruct
type TestStruct struct {
	Val  int
	Name string
}

// Reset сбрасывает структуру к нулевым значениям
func (t *TestStruct) Reset() {
	t.Val = 0
	t.Name = ""
}

// TestPool_GetPut - проверяет базовую логику Get/Put и вызов Reset()
func TestPool_GetPut(t *testing.T) {
	// 1. Создаем пул
	// T - это *TestStruct, который реализует Resetter
	// newItem() создает новый *TestStruct
	pool := New(func() *TestStruct {
		return new(TestStruct)
	})

	// 2. Get() из пустого пула (вызывает New)
	obj1 := pool.Get()
	assert.NotNil(t, obj1)
	assert.Equal(t, 0, obj1.Val)
	assert.Equal(t, "", obj1.Name)

	// 3. Модифицируем объект
	obj1.Val = 100
	obj1.Name = "Hello"

	// 4. Put() - возвращаем в пул. Reset() должен быть вызван немедленно
	pool.Put(obj1)

	// 5. Get() - получаем тот же самый (сброшенный) объект
	obj2 := pool.Get()
	assert.NotNil(t, obj2)

	// Проверяем, что Reset() был вызван
	assert.Equal(t, 0, obj2.Val, "Поле Val должно было сброситься")
	assert.Equal(t, "", obj2.Name, "Поле Name должно было сброситься")

	// Проверяем, что это тот же указатель (тот же объект в памяти)
	assert.Same(t, obj1, obj2, "Объект должен быть переиспользован из пула")
}

// TestPool_New - проверяет, что newItem вызывается только, когда пул пуст
func TestPool_New(t *testing.T) {
	newCounter := 0
	pool := New(func() *TestStruct {
		newCounter++
		return &TestStruct{Val: -1} // Помечаем, что он новый
	})

	// 1. Get() - newCounter должен стать 1
	obj1 := pool.Get()
	assert.Equal(t, 1, newCounter)
	assert.Equal(t, -1, obj1.Val) // Убеждаемся, что это новый

	obj1.Val = 100
	pool.Put(obj1)
	assert.Equal(t, 1, newCounter) // Put() не вызывает New

	// 2. Get() - newCounter не должен измениться (получаем из пула)
	obj2 := pool.Get()
	assert.Equal(t, 1, newCounter)
	assert.Equal(t, 0, obj2.Val) // Сброшенный (0), а не новый (-1)
	assert.Same(t, obj1, obj2)

	// 3. Get() - пул снова пуст, newCounter должен стать 2
	obj3 := pool.Get()
	assert.Equal(t, 2, newCounter)
	assert.Equal(t, -1, obj3.Val) // Снова новый
}

// TestPool_Concurrency - проверяет потокобезопасность
func TestPool_Concurrency(t *testing.T) {
	pool := New(func() *TestStruct {
		return new(TestStruct)
	})

	var wg sync.WaitGroup
	numGoroutines := 100
	numOps := 1000

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				obj := pool.Get()
				assert.Equal(t, 0, obj.Val) // Должен быть сброшен

				// "Работаем" с объектом
				obj.Val = id
				obj.Name = "goroutine-" + strconv.Itoa(id)

				pool.Put(obj)
			}
		}(i)
	}

	wg.Wait()
	// Тест пройден, если не было 'race' (запускать с -race) и assert не упали
}
