package cli


import (
	"fmt"
	"strconv"
	"os"
	"strings"
	"errors"

	"text/tabwriter"

	"github.com/Tinch334/Computer-one-v2/co"
)


func breakpointHandler(ctrl *interpreterControl, cfg *interpreterConfig, args []string) {
    switch args[0] {
    case BREAKPOINT_SET:
    	if len(args) != 2 {
            printErrorMsg("breakpoint")
            return
        }

        //Get line number and check for errors.
        err, addr := convValidateMemoryAddr(args[1])
        if err != nil {
        	return
        }

        ctrl.AddBreakpoint(addr)
        fmt.Printf("Breakpoint added at address 0x%X", addr)

    case BREAKPOINT_LIST:
    	br := ctrl.GetBreakpoints()

    	if len(br) == 0 {
    		fmt.Printf("No breakpoints set")
    	} else {
    		fmt.Printf("Breakpoints set at addresses: %s", strings.Join(sliceMap(br, uint16ToHexStr()), ", "))
    	}

    case BREAKPOINT_DELETE:
    	//Get line number and check for errors.
        err, addr := convValidateMemoryAddr(args[1])
        if err != nil {
        	return
        }

        if ctrl.HasBreakpoint(addr) {
        	ctrl.DeleteBreakpoint(addr)
        	fmt.Printf("Breakpoint successfully deleted")
        } else {
        	fmt.Printf("Breakpoint not found")
        }

    default:
        printErrorMsg("breakpoint")
    }
}

func printHelp() {
    //Use tab-writer for easy alignment.
    tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)

    fmt.Fprintln(tw, "Available commands:\n")

    type cmd struct {
        name     string
        short    string
        desc     string
        options []string
    }

    cmds := []cmd{
        {name: STEP, short: STEP_SHORT, desc: "Perform one execution step"},
        {name: CONTINUE, short: CONTINUE_SHORT, desc: "Continue execution until a breakpoint or program end"},
        {
            name:  BREAKPOINT,
            short: BREAKPOINT_SHORT,
            desc:  "Breakpoint handler, options:",
            options: []string{
                fmt.Sprintf("%s <address>\tSet breakpoint at <address>", BREAKPOINT_SET),
                fmt.Sprintf("%s\tList all breakpoints", BREAKPOINT_LIST),
                fmt.Sprintf("%s <address>\tDelete the breakpoint at <address>, if it exists", BREAKPOINT_DELETE),
            },
        },
        {name: EXIT, short: EXIT_SHORT, desc: "Exit interpreter"},
        {name: HELP, short: HELP_SHORT, desc: "Display this help message"},
        {
            name: CONFIGURE,
            short: CONFIGURE_SHORT,
            desc: "Configure the interpreter",
            options: []string{
                fmt.Sprintf("%s <lower> <upper>\tSets the bounds determining which memory cells are printed", CONFIGURE_MEMORY_LIMITS),
            },
        },
    }

    for _, c := range cmds {
        fmt.Fprintf(tw, "%s\t- %s\t| %s\n", c.name, c.short, c.desc)
        //Check if there are any options.
        if len(c.options) > 0 {
            for _, ex := range c.options {
                fmt.Fprintf(tw, "\t\t\t%s\n", ex)
            }
        }
    }

    //Ensure no elements left in buffer.
    _ = tw.Flush()
}

func printErrorMsg(functionName string)  {
    fmt.Fprintf(os.Stderr, "Invalid usage for \"%s\", see help", functionName)
}

//A function that converts the given uint16 number to a string with it's hex representation.
func uint16ToHexStr() func(uint16) string {
	return func(num uint16) string {
		return fmt.Sprintf("0x%X", num)
	}	
}

//Takes a string and if possible converts it to a uint16 number, otherwise returns an error.
func convValidateMemoryAddr(addr string) (error, uint16) {
	//Try base 10 conversion first.
	num, err := strconv.ParseUint(addr, 10, 16)
    if err != nil {
    	//If base 10 conversion failed try base 16 conversion, we force base 16 numbers to begin with "0x".
    	if strings.HasPrefix(addr, "0x") {
    		num, err = strconv.ParseUint(addr[2:], 16, 16)
    	} else {
    		printErrorMsg("breakpoint")
    		return errors.New("Invalid memory address"), 0
    	}

    	if err != nil {
    		printErrorMsg("breakpoint")
    		return errors.New("Invalid memory address"), 0
    	}
    }

    if num >= co.MemorySize {
    	printErrorMsg("breakpoint")
    	return errors.New("Invalid memory address"), 0
    }

    return nil, uint16(num)
}