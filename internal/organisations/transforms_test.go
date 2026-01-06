package organisations

import (
	"testing"
	"time"

	"github.com/code-gorilla-au/odize"
)

func TestTransforms(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("toTime should convert unix timestamp correctly", func(t *testing.T) {
			ts := int64(1609459200) // 2021-01-01 00:00:00 UTC
			expected := time.Unix(ts, 0).UTC()
			actual := toTime(ts)
			odize.AssertTrue(t, expected.Equal(actual))
		}).
		Run()

	odize.AssertNoError(t, err)
}
