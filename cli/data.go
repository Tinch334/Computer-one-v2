package cli


import (
	"slices"

	"github.com/fatih/color"
)


type interpreterControl struct {
	running bool

	breakpoints []uint16

	step bool
	cont bool
}

type interpreterConfig struct {
	memoryLimitL, memoryLimitH uint16

	highlightPC bool
	highlightPCColour *color.Color

	exitOnError bool
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
	BREAKPOINT_DELETE_ALL = "da"

	EXIT = "exit"
	EXIT_SHORT = "e"

	HELP = "help"
	HELP_SHORT = "h"

	CONFIGURE = "configure"
	CONFIGURE_SHORT = "cfg"

	CONFIGURE_MEMORY_LIMITS = "ml"
)


func (c *interpreterControl) AddBreakpoint(pos uint16) {
	//Avoid duplicate breakpoints.
	if !slices.Contains(c.breakpoints, pos) {
		c.breakpoints = append(c.breakpoints, pos)
	}
}

func (c *interpreterControl) HasBreakpoint(pos uint16) bool {
	return slices.Contains(c.breakpoints, pos)
}

func (c *interpreterControl) DeleteBreakpoint(pos uint16) {
	del := func (e uint16) bool {
		return e == pos
	}

	c.breakpoints = slices.DeleteFunc(c.breakpoints, del)
}

func (c *interpreterControl) ClearBreakpoints() {
	c.breakpoints = make([]uint16, 0)
}

func (c* interpreterControl) GetBreakpoints() []uint16 {
	return c.breakpoints
}


func (cfg *interpreterConfig) SetMemoryLimits(l, h uint16) {
	cfg.memoryLimitL = l
	cfg.memoryLimitH = h
}