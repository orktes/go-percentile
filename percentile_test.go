package percentile

import (
	"testing"
	"time"
)

func TestCalculatorForP50(t *testing.T) {
	c := NewCalculator(time.Millisecond, time.Duration(0), time.Millisecond*1000)
	c.Add(time.Millisecond)

	p50 := c.Percentile(0.5)
	if p50 != time.Millisecond {
		t.Error("wrong value returned", p50)
	}

	c = NewCalculator(time.Millisecond, time.Duration(0), time.Millisecond*1000)
	c.Add(1 * time.Millisecond)
	c.Add(2 * time.Millisecond)
	c.Add(3 * time.Millisecond)

	if c.Percentile(0.5) != 2*time.Millisecond {
		t.Error("wrong p50 returned", c.Percentile(0.5))
	}

	c.Add(4 * time.Millisecond)

	if c.Percentile(0.5) != 2500*time.Microsecond {
		t.Error("wrong p50 returned", c.Percentile(0.5))
	}

	c.Add(4 * time.Millisecond)

	if c.Percentile(0.5) != 3000*time.Microsecond {
		t.Error("wrong p50 returned", c.Percentile(0.5))
	}
}

func TestMultiCalculatorForP50(t *testing.T) {
	c1 := NewCalculator(time.Millisecond, time.Duration(0), time.Millisecond*1000)
	c2 := NewCalculator(time.Millisecond, time.Duration(0), time.Millisecond*1000)
	c3 := NewCalculator(time.Millisecond, time.Duration(0), time.Millisecond*1000)

	c1.Add(1 * time.Millisecond)
	c2.Add(2 * time.Millisecond)
	c3.Add(3 * time.Millisecond)

	d := MultiPercentile(0.5, c1, c2, c3)
	if d != 2*time.Millisecond {
		t.Error("Wrong p50 returned")
	}
}
