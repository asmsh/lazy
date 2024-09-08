package lazy_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/asmsh/lazy"
)

func TestNew(t *testing.T) {
	t.Run("when no error", func(t *testing.T) {
		value := lazy.NewValue(func() (int, error) {
			return 123, nil
		})

		got := value.Val()
		if got != 123 {
			t.Errorf("got %d, want %d", got, 123)
		}

		err := value.Err()
		if err != nil {
			t.Errorf("got %v, want %v", err, nil)
		}
	})

	t.Run("when error", func(t *testing.T) {
		value := lazy.NewValue(func() (int, error) {
			return 123, errors.New("error")
		})

		got := value.Val()
		if got != 123 {
			t.Errorf("got %d, want %d", got, 123)
		}

		err := value.Err()
		if err == nil {
			t.Errorf("got %v, want %v", err, errors.New("error"))
		}
	})

	t.Run("not loaded then loaded", func(t *testing.T) {
		value := lazy.NewValue(func() (int, error) {
			return 123, errors.New("error")
		})

		gotBefore := value.IsLoaded()
		if gotBefore != false {
			t.Errorf("got %t, want %t", gotBefore, false)
		}

		value.Val() // or value.Err()

		gotAfter := value.IsLoaded()
		if gotAfter != true {
			t.Errorf("got %t, want %t", gotAfter, true)
		}
	})

	t.Run("concurrent reading", func(t *testing.T) {
		value := lazy.NewValue(func() (int, error) {
			return 123, errors.New("error")
		})

		wg := sync.WaitGroup{}
		wg.Add(3) // 3 concurrent readers

		// concurrent reader 1
		go func() {
			defer wg.Done()
			got := value.Val()
			if got != 123 {
				t.Errorf("got %d, want %d", got, 123)
			}
		}()

		// concurrent reader 2
		go func() {
			defer wg.Done()
			got := value.Val()
			if got != 123 {
				t.Errorf("got %d, want %d", got, 123)
			}
		}()

		// concurrent error handler
		go func() {
			defer wg.Done()
			err := value.Err()
			if err == nil {
				t.Errorf("got %v, want %v", err, errors.New("error"))
			}
		}()

		wg.Wait()
	})
}
