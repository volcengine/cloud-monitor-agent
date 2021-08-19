package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mackerelio/go-osstat/uptime"
	"github.com/pkg/errors"
)

// GetUptimeInSeconds return the uptime.
func GetUptimeInSeconds() (int64, error) {
	osUptime, err := uptime.Get()
	if err != nil {
		return -1, errors.Wrap(err, "exec uptime.Get failed")
	}

	return int64(osUptime / time.Second), nil
}

// Float64From32Bits the value of f maybe equal zero, so this time we return the zero directly.
func Float64From32Bits(f float64) float64 {
	if f < 0 {
		return 0
	}
	return f
}

// Keep2Decimal return the float which have two decimal.
func Keep2Decimal(number float64) float64 {
	limitedString := fmt.Sprintf("%.2f", number)
	limitedNumber, _ := strconv.ParseFloat(limitedString, 64)
	return limitedNumber
}
