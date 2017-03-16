package expr

import (
	"fmt"
	"strings"
)

func toTrimZero(f float64) string {
	spow := fmt.Sprintf("%f", f)
	spow = strings.TrimRight(spow, ".0")
	return spow
}
