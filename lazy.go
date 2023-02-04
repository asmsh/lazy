package lazy

import "sync"

// Value is a generic container offering lazy-loading for values of any type.
// Safe for concurrent usage.
type Value[T any] interface {
	// Val returns the lazily-loaded value contained within this Value.
	// It will load the value if it wasn't loaded yet.
	// All callers will block until the value is loaded.
	Val() T

	// Err returns the error that occurred while loading the value.
	// It will load the value if it wasn't loaded yet.
	// All callers will block until the value is loaded.
	Err() error

	private()
}

// NewValue creates a new Value, which will be lazily-loaded from the loader
// function provided, init, upon first call to any of its methods.
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

func (lv *lazyValue[T]) private() {}

func (lv *lazyValue[T]) Val() T {
	lv.load()
	return lv.val
}

func (lv *lazyValue[T]) Err() error {
	lv.load()
	return lv.err
}
