package cli


import (
	"slices"

	"github.com/fatih/color"
)


type interpreterConfig struct {
	memoryLimitL, memoryLimitH uint16

	highlightPC bool
	highlightPCColour *color.Color
}

type interpreterControl struct {
	running bool

	breakpoints []uint16

	step bool
	cont bool
}

const (
	STEP = "step"
	STEP_SHORT = "s"

	CONTINUE = "continue"
	CONTINUE_SHORT = "c"

	BREAKPOINT = "breakpoint"
	BREAKPOINT_SHORT = "br"

	BREAKPOINT_SET = "s"
	BREAKPOINT_LIST = "l"
	BREAKPOINT_DELETE = "d"

	EXIT = "exit"
	EXIT_SHORT = "e"

	HELP = "help"
	HELP_SHORT = "h"
)


func (c *interpreterControl) AddBreakpoint(pos uint16) {
	// avoid duplicate breakpoints
	if !slices.Contains(c.breakpoints, pos) {
		c.breakpoints = append(c.breakpoints, pos)
	}
}

func (c *interpreterControl) HasBreakpoint(pos uint16) bool {
	return slices.Contains(c.breakpoints, pos)
}

func (cfg *interpreterConfig) SetMemoryLimits(l, h uint16) {
	cfg.memoryLimitL = l
	cfg.memoryLimitH = h
}