package scheduling

import (
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
	t.Run("should call the given function, at the given interval", func(t *testing.T) {
		defer revertStubs()

		now := time.Now()
		next := now.Add(1 * time.Millisecond)

		timeNow = timeNowTest
		timeNowTestTime = now

		parseExpr = parseExprTest
		parseExprTestExpr = &testExpression{
			next: next,
		}

		errs := make(chan error)
		quit := make(chan int)
		call := make(chan string)

		go ScheduleFunc(quit, errs, "* * * * * * *", func() error {
			call <- "called"
			quit <- 0

			return nil
		})

		select {
		case res := <-call:
			assert.Equal(t, "called", res)
		case err := <-errs:
			t.Fatal(err)
		}
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
