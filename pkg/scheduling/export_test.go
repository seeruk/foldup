package scheduling

import "time"

var (
	parseExprTestExpr expression
	parseExprTestErr  error
	timeNowTestTime   time.Time
)

type testExpression struct {
	next time.Time
}

func (e *testExpression) Next(time.Time) time.Time {
	return e.next
}

func parseExprTest(expr string) (expression, error) {
	return parseExprTestExpr, parseExprTestErr
}

func timeNowTest() time.Time {
	return timeNowTestTime
}

func revertStubs() {
	parseExpr = parseCronExpr
	timeNow = time.Now
}
