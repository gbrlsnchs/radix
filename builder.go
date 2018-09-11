package radix

import (
	"strings"

	"github.com/fatih/color"
)

type builder struct {
	*strings.Builder
	colors [4]*color.Color
	debug  bool
}
