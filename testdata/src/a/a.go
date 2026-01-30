package a

import "time"

// TimeAlias is an alias for time.Time
type TimeAlias = time.Time

// CustomTime is a custom type with Add method (should not be detected)
type CustomTime struct{}

func (c CustomTime) Add(d time.Duration) CustomTime {
	return c
}

func (c CustomTime) Sub(other CustomTime) time.Duration {
	return 0
}

func example() {
	t := time.Now()
	d := time.Hour

	// Direct value calls - should be detected
	_ = t.Add(d)     // want "use of time.Time.Add is not allowed"
	_ = t.AddDate(1, 0, 0) // want "use of time.Time.AddDate is not allowed"
	_ = t.Sub(t)     // want "use of time.Time.Sub is not allowed"
	_ = t.Truncate(d) // want "use of time.Time.Truncate is not allowed"
	_ = t.Round(d)   // want "use of time.Time.Round is not allowed"

	// Pointer calls - should be detected
	pt := &t
	_ = pt.Add(d)    // want "use of time.Time.Add is not allowed"
	_ = pt.Sub(t)    // want "use of time.Time.Sub is not allowed"

	// Type alias - should be detected
	var ta TimeAlias = time.Now()
	_ = ta.Add(d)    // want "use of time.Time.Add is not allowed"
	_ = ta.Sub(t)    // want "use of time.Time.Sub is not allowed"

	// Custom type with same method names - should NOT be detected
	ct := CustomTime{}
	_ = ct.Add(d)    // OK - not time.Time
	_ = ct.Sub(ct)   // OK - not time.Time

	// Allowed methods - should NOT be detected
	_ = t.Year()
	_ = t.Month()
	_ = t.Day()
	_ = t.Format(time.RFC3339)
	_ = t.Unix()
	_ = t.UnixNano()
	_ = t.Before(t)
	_ = t.After(t)
	_ = t.Equal(t)
	_ = t.IsZero()
	_ = t.Location()
	_, _ = t.Zone()
}

func chainedCalls() {
	t := time.Now()
	d := time.Hour

	// Chained calls - should be detected
	_ = t.Add(d).Add(d) // want "use of time.Time.Add is not allowed" "use of time.Time.Add is not allowed"
}

func functionReturningTime() time.Time {
	return time.Now()
}

func callOnReturnValue() {
	d := time.Hour

	// Call on function return value - should be detected
	_ = functionReturningTime().Add(d) // want "use of time.Time.Add is not allowed"
	_ = time.Now().Add(d)              // want "use of time.Time.Add is not allowed"
}
