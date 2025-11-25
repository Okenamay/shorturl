package pool

import "sync"

// Resetter определяет интерфейс для объектов, которые можно "сбросить", что
// позволяет переиспользовать их в пуле
type Resetter interface {
	Reset()
}

// Pool - это потокобезопасный пул объектов, использующий дженерики, он хранит
// объекты типа T, которые должны реализовывать интерфейс Resetter
type Pool[T Resetter] struct {
	p *sync.Pool
}

// New создает новый пул. Он требует функцию `newItem`, которая будет
// вызываться для создания нового объекта T, когда пул пуст
func New[T Resetter](newItem func() T) *Pool[T] {
	return &Pool[T]{
		p: &sync.Pool{
			// sync.Pool требует, чтобы New возвращал any (он же interface{})
			New: func() any {
				return newItem()
			},
		},
	}
}

// Get извлекает объект T из пула. Если пул пуст, будет создан новый объект с
// помощью функции newItem
func (p *Pool[T]) Get() T {
	// Мы безопасно приводим тип any обратно к T
	return p.p.Get().(T)
}

// Put возвращает объект T в пул. Перед помещением в пул, у объекта
// автоматически вызывается метод Reset()
func (p *Pool[T]) Put(obj T) {
	// Это главная логика: T реализует Resetter, поэтому у него есть метод
	// Reset()
	obj.Reset()
	p.p.Put(obj)
}
