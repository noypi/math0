package expr

import (
	"fmt"
	"strings"
)

func toTrimZero(f float64) string {
	spow := fmt.Sprintf("%f", f)
	spow = strings.TrimRight(spow, "0")
	if '.' == spow[len(spow)-1] {
		spow = spow[:len(spow)-1]
	}
	return spow
}
