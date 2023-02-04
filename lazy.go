package lazy

import "sync"

// Value is a generic Lazy loader.
type Value[T any] struct {
	init func() (T, error)
	once sync.Once
	val  T
	err  error
}

// NewValue creates a new Value.
func NewValue[T any](init func() (T, error)) Value[T] {
	return Value[T]{init: init}
}

func (c *Value[T]) load() {
	c.once.Do(func() {
		c.val, c.err = c.init()
		// release init, so the GC can collect it.
		c.init = nil
	})
}

// Val loads if the value is not loaded yet,
// and returns the value.
func (c *Value[T]) Val() T {
	c.load()
	return c.val
}

func (c *Value[T]) Err() error {
	c.load()
	return c.err
}
