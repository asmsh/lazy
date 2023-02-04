package lazy

import "sync"

// Loader is a generic Lazy loader.
type Loader[T any] struct {
	onc      sync.Once
	val      T
	err      error
	supplier func() (T, error)
}

// New creates a new Loader.
func New[T any](supplier func() (T, error)) Loader[T] {
	return Loader[T]{
		onc:      sync.Once{},
		supplier: supplier,
	}
}

func (c *Loader[T]) load() {
	c.onc.Do(func() {
		val, err := c.supplier()
		c.val = val
		c.err = err
		// release the supplier, so the GC can collect it.
		c.supplier = nil
	})
}

// Value loads if the value is not loaded yet,
// and returns the value.
func (c *Loader[T]) Value() T {
	c.load()
	return c.val
}

func (c *Loader[T]) Error() error {
	c.load()
	return c.err
}
