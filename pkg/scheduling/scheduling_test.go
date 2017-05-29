package scheduling

import (
	"errors"
	"testing"
	"time"

	"github.com/SeerUK/assert"
)

func TestParseCronExpr(t *testing.T) {
	t.Run("should error if given an invalid expression", func(t *testing.T) {
		_, err := parseCronExpr("Hello, World!")

		assert.NotOK(t, err)
	})

	t.Run("should return a parsed expression that yields the next time", func(t *testing.T) {
		expr, err := parseCronExpr("* * * * * *")
		assert.OK(t, err)

		layout := time.RFC822

		// Remove seconds from time, since the minimum we accuracy of a cron
		// expression is minutes.
		now, err := time.Parse(layout, time.Now().Format(layout))
		assert.OK(t, err)

		expected := now.Add(1 * time.Minute)

		assert.Equal(t, expected, expr.Next(now))
	})
}

func TestScheduleFunc(t *testing.T) {
	t.Run("should call the given function, at the given interval, and quit", func(t *testing.T) {
		defer revertStubs()

		now := time.Now()
		next := now.Add(1 * time.Millisecond)

		timeNow = timeNowTest
		timeNowTestTime = now

		parseExpr = parseExprTest
		parseExprTestExpr = &testExpression{
			next: next,
		}

		call := make(chan bool, 1)
		done := make(chan int, 1)

		err := ScheduleFunc(done, "* * * * * * *", func() error {
			call <- true
			done <- 1

			return nil
		})

		assert.OK(t, err)
		assert.Equal(t, true, <-call)
	})

	t.Run("should error if an invalid cron expression is passed", func(t *testing.T) {
		done := make(chan int, 1)

		err := ScheduleFunc(done, "hello world", func() error {
			return nil
		})

		assert.NotOK(t, err)
	})

	t.Run("should error if an error is returned in the callback", func(t *testing.T) {
		defer revertStubs()

		now := time.Now()
		next := now.Add(1 * time.Millisecond)

		timeNow = timeNowTest
		timeNowTestTime = now

		parseExpr = parseExprTest
		parseExprTestExpr = &testExpression{
			next: next,
		}

		done := make(chan int, 1)

		err := ScheduleFunc(done, "* * * * * * *", func() error {
			return errors.New("This is an error")
		})

		assert.NotOK(t, err)
	})
}
