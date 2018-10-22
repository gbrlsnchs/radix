package radix

import (
	"strings"

	"github.com/gbrlsnchs/color"
)

type builder struct {
	*strings.Builder
	colors [4]color.Color
	debug  bool
}
