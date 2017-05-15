package scheduling

import (
	"testing"
	"time"

	"github.com/SeerUK/assert"
)

func TestParseCronExpr(t *testing.T) {
	t.Run("should error if given an invalid expression", func(t *testing.T) {
		t.Skip()
	})

	t.Run("should return a parsed expression that yields the next time", func(t *testing.T) {
		t.Skip()
	})
}

func TestScheduleFunc(t *testing.T) {
	t.Run("should call the given function, at the given interval", func(t *testing.T) {
		quit := make(chan int)
		called := false

		go ScheduleFunc(quit, "* * * * * * *", func() error {
			called = true

			return nil
		})

		time.Sleep(1 * time.Second)

		quit <- 0

		assert.True(t, called, "Expected callback to have been called.")
	})

	t.Run("should quit when asked, as quickly as possible", func(t *testing.T) {
		t.Skip()
	})

	t.Run("should error if an invalid cron expression is passed", func(t *testing.T) {
		t.Skip()
	})

	t.Run("should error if an error is returned in the callback", func(t *testing.T) {
		t.Skip()
	})
}
