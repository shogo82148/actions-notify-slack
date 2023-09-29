package database

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type attrConverter struct {
	err error
}

func (c *attrConverter) convertString(attr types.AttributeValue) string {
	if c.err != nil {
		return ""
	}
	s, ok := attr.(*types.AttributeValueMemberS)
	if !ok {
		c.err = fmt.Errorf("database: %T cannot convert into string", attr)
		return ""
	}
	return s.Value
}

func (c *attrConverter) convertNumber(attr types.AttributeValue) float64 {
	if c.err != nil {
		return 0
	}
	n, ok := attr.(*types.AttributeValueMemberN)
	if !ok {
		c.err = fmt.Errorf("database: %T cannot convert into number", attr)
		return 0
	}
	f, err := strconv.ParseFloat(n.Value, 64)
	if err != nil {
		c.err = err
		return 0
	}
	return f
}

func timeToUnixTime(t time.Time) float64 {
	unix := t.Unix()
	nano := t.Nanosecond()
	return math.FMA(float64(nano), 1e-9, float64(unix))
}

func unixTimeToTime(unixTime float64) time.Time {
	i, f := math.Modf(unixTime)
	return time.Unix(int64(i), int64(math.RoundToEven(f*1e9)))
}
