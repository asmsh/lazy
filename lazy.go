package lazy

import "sync"

// Value is a generic Lazy loader.
type Value[T any] interface {
	// Val loads if the value is not loaded yet,
	// and returns the value.
	Val() T

	Err() error
}

// NewValue creates a new Value.
func NewValue[T any](init func() (T, error)) Value[T] {
	return &lazyValue[T]{init: init}
}

type lazyValue[T any] struct {
	init func() (T, error)
	once sync.Once
	val  T
	err  error
}

func (lv *lazyValue[T]) load() {
	lv.once.Do(func() {
		lv.val, lv.err = lv.init()
		// release init, so the GC can collect it.
		lv.init = nil
	})
}

func (lv *lazyValue[T]) Val() T {
	lv.load()
	return lv.val
}

func (lv *lazyValue[T]) Err() error {
	lv.load()
	return lv.err
}
