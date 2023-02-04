package lazy_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/asmsh/lazy"
)

func TestNew(t *testing.T) {
	t.Run("when no error", func(t *testing.T) {
		loader := lazy.New(func() (int, error) {
			return 123, nil
		})

		if loader.Loaded() {
			t.Errorf("got %t, want %t", loader.Loaded(), false)
		}

		got := loader.Value()
		if got != 123 {
			t.Errorf("got %d, want %d", got, 123)
		}

		if !loader.Loaded() {
			t.Errorf("got %t, want %t", loader.Loaded(), true)
		}

		err := loader.Error()
		if err != nil {
			t.Errorf("got %v, want %v", err, nil)
		}
	})

	t.Run("when error", func(t *testing.T) {
		loader := lazy.New(func() (int, error) {
			return 123, errors.New("error")
		})

		if loader.Loaded() {
			t.Errorf("got %t, want %t", loader.Loaded(), false)
		}

		got := loader.Value()
		if got != 123 {
			t.Errorf("got %d, want %d", got, 123)
		}

		if !loader.Loaded() {
			t.Errorf("got %t, want %t", loader.Loaded(), true)
		}

		err := loader.Error()
		if err == nil {
			t.Errorf("got %v, want %v", err, errors.New("error"))
		}
	})

	t.Run("concurrent reading", func(t *testing.T) {
		loader := lazy.New(func() (int, error) {
			return 123, errors.New("error")
		})

		wg := sync.WaitGroup{}
		wg.Add(3) // 3 concurrent readers

		// concurrent reader 1
		go func() {
			defer wg.Done()
			got := loader.Value()
			if got != 123 {
				t.Errorf("got %d, want %d", got, 123)
			}
		}()

		// concurrent reader 2
		go func() {
			defer wg.Done()
			got := loader.Value()
			if got != 123 {
				t.Errorf("got %d, want %d", got, 123)
			}
		}()

		// concurrent error handler
		go func() {
			defer wg.Done()
			err := loader.Error()
			if err == nil {
				t.Errorf("got %v, want %v", err, errors.New("error"))
			}
		}()

		wg.Wait()
	})
}
