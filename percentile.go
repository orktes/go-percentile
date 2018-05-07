package percentile

import (
	"errors"
	"sync"
	"time"
)

// ErrDurationBelowMinLimit is returned when given duration is below min limit
var ErrDurationBelowMinLimit = errors.New("value below limit")

// ErrDurationAboveMinLimit is returned when given duration is below min limit
var ErrDurationAboveMinLimit = errors.New("value above limit")

// Calculator stores values and calculates percentiles for given values
type Calculator struct {
	bucketSize time.Duration
	maxValue   time.Duration
	minValue   time.Duration

	size   uint64
	values []uint64
	count  uint64
	mutex  sync.RWMutex
}

// NewCalculator returns a new percentile calculator for
// given bucket size, min val & max val
func NewCalculator(bucketSize time.Duration, minValue time.Duration, maxValue time.Duration) *Calculator {
	valueRange := maxValue - minValue
	size := valueRange / bucketSize

	return &Calculator{
		minValue:   minValue,
		maxValue:   maxValue,
		bucketSize: bucketSize,
		size:       uint64(size),
		values:     make([]uint64, size, size),
	}
}

// Add add value to calculator
func (c *Calculator) Add(d time.Duration) error {
	if d < c.minValue {
		return ErrDurationBelowMinLimit
	} else if d > c.maxValue {
		return ErrDurationAboveMinLimit
	}

	indx := (d - c.minValue) / c.bucketSize

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.count++
	c.values[int(indx)]++
	return nil
}

// Percentile returns a percentile between 0.0 - 1.0 from the calucalor
func (c *Calculator) Percentile(percentile float64) time.Duration {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	deciredIndex := float64(c.count) * percentile
	vindex := uint64(0)
	even := deciredIndex == float64(int64(deciredIndex))

	var lowBucket time.Duration

	for i, count := range c.values {
		vindex = vindex + count
		if deciredIndex <= float64(vindex) {
			val := c.minValue + time.Duration(i)*c.bucketSize

			if lowBucket > 0 {
				return (lowBucket + val) / 2
			}

			if even {
				lowBucket = val
				continue
			}
			return val
		}
	}

	return c.minValue + time.Duration(c.values[len(c.values)-1])*c.bucketSize
}

// MultiPercentile returns percentile from multiple similarly configured calculators
func MultiPercentile(percentile float64, calcs ...*Calculator) time.Duration {
	totalCount := uint64(0)
	for _, c := range calcs {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		totalCount += c.count
	}

	c := calcs[0]

	deciredIndex := float64(totalCount) * percentile
	vindex := uint64(0)
	even := deciredIndex == float64(int64(deciredIndex))

	var lowBucket time.Duration

	for duration := range c.values {
		count := uint64(0)
		for _, calc := range calcs {
			count += calc.values[duration]
		}

		vindex = vindex + count
		if deciredIndex <= float64(vindex) {
			val := c.minValue + time.Duration(duration)*c.bucketSize

			if lowBucket > 0 {
				return (lowBucket + val) / 2
			}

			if even {
				lowBucket = val
				continue
			}
			return val
		}
	}

	return c.minValue + time.Duration(c.values[len(c.values)-1])*c.bucketSize
}
