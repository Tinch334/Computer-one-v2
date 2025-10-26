package cli

import (
    "fmt"
    "strings"
    "bufio"
    "os"

    "text/tabwriter"

    "github.com/Tinch334/Computer-one-v2/co"
    "github.com/fatih/color"
)

func RunCli() {
    ci := co.NewComputerInfo()

    var memLoad []uint16
    memLoad = append(
        memLoad,
        0b0000100000000011, //LDA r2 3
        0b0010000000000010, //LDR r0 r2
        0b1110000000000000, //HLT
        0b0000000000001110,
        0b0000000000000000,
        0b0000000000000000,
        0b0000000000000000,
        0b0000000000000000,
        0b0000000000000000,
        0b0000000000000000,
        0b0000000000000000,
        0b0000000000000000,
        0b0000000000000000,
        0b0000000000000000,
        0b0000000000000011,
    )

    ci.SetMemory(0, memLoad)

    reader := bufio.NewReader(os.Stdin)

    control := interpreterControl{
        running: true,
        step: false,
        cont: false,
    }

    config := interpreterConfig {
        memoryLimitL: 0,
        memoryLimitH: 40,

        highlightPC: true,
        highlightPCColour: color.New(color.FgBlue),

        exitOnError: false,
    }

    run(ci, reader, &control, &config)
}

func run(ci *co.ComputerInfo, reader *bufio.Reader, ctrl *interpreterControl, cfg *interpreterConfig) {
    printNext := true

    //Runs the interpreting loop as long as the program runs.
    for ctrl.running {
        //Print info.
        if printNext {
            printRegs(ci)
            printMemory(ci, cfg)

            printNext = false
        }

        ctrl.step = false

        //Check if we should continue or check for a step.
        if ctrl.cont {
            //Step the program until a breakpoint is reached
            if ctrl.HasBreakpoint(ci.GetRegisters().PC) {
                ctrl.cont = false
            } else {
                ctrl.step = true
            }
            
        } else {
            processInput(reader, ctrl, cfg)
        }

        //Step program.
        if ctrl.step {
            err, run := ci.Step()
            if !run {
                ctrl.running = false
                fmt.Printf("Program halted\n")
            }

            if err != nil {
                if cfg.exitOnError {
                    ctrl.running = false
                } else {
                    fmt.Printf("An error occurred during execution: %s\n", err)
                }
            }

            printNext = true
        }

        fmt.Printf("\n")
    }
}

func processInput(reader *bufio.Reader, ctrl *interpreterControl, cfg *interpreterConfig) {
    fmt.Printf(">")
    line, err := reader.ReadString('\n')

    if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }
    
    contents := strings.Fields(line)
    command, arguments := contents[0], contents[1:]

    switch command {
    case STEP:
        fallthrough
    case STEP_SHORT:
        ctrl.step = true

    case CONTINUE:
        fallthrough
    case CONTINUE_SHORT:
        ctrl.cont = true

    case BREAKPOINT:
        fallthrough
    case BREAKPOINT_SHORT:
        if len(arguments) == 0 {
            printErrorMsg("breakpoint")
        }

        switch arguments[0] {
        case BREAKPOINT_SET:
            /* code */
        default:
            printErrorMsg("breakpoint")
            return
        }

    case EXIT:
        fallthrough
    case EXIT_SHORT:
        ctrl.running = false

    case HELP:
        fallthrough
    case HELP_SHORT:
        printHelp()

    case CONFIGURE:
        fallthrough
    case CONFIGURE_SHORT:
        printHelp()

    default:
        fmt.Println("Unknown command, use \"h\" for help")
        return
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
                fmt.Sprintf("%s <line>\tSet breakpoint at <line>", BREAKPOINT_SET),
                fmt.Sprintf("%s\tList all breakpoints", BREAKPOINT_LIST),
                fmt.Sprintf("%s <line>\tDelete the breakpoint at <line>, if it exists", BREAKPOINT_DELETE),
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

    _ = tw.Flush()
}

func printErrorMsg(functionName string) {
    fmt.Printf("Invalid usage for \"%s\", see help\n", functionName)
}


/*
    DISPLAY FUNCTIONS
*/
//Returns a "1" if the given boolean is true, "0" otherwise.
func btoi(b bool) string {
    if b {
        return "1"
    }
    return "0"
}

//Prints all registers.
func printRegs(ci *co.ComputerInfo) {
    regs := ci.GetRegisters()
    flags := ci.GetFlags()
    flagsStr := btoi(flags.N) + btoi(flags.P) + btoi(flags.Z)

    fmt.Printf("PC: 0x%04x | NPZ: %s | R0: 0x%04x R1: 0x%04x R2: 0x%04x R3:0x%04x R4:0x%04x R5:0x%04x RR: 0x%04x\n",
        regs.PC, flagsStr, regs.R0, regs.R1, regs.R2, regs.R3, regs.R4, regs.R5, regs.RR)
}

//Prints memory contents.
func printMemory(ci *co.ComputerInfo, cfg *interpreterConfig) {
    start := cfg.memoryLimitL
    end := cfg.memoryLimitH

    pc := int(ci.GetRegisters().PC)

    valuesPerRow := 8
    const hexAddrWidth = 4

    //Get appropriate memory cells to print.
    _, mem := ci.GetMemory(start, end)

    for i, elem := range mem {
        //Print memory cell number along with spacing.
        if (i % valuesPerRow) == 0 {
            if i != 0 {
                fmt.Printf("\n")
            }
            fmt.Printf("0x%0*x : ", hexAddrWidth, (start + uint16(i)))
        }
        //Print element.
        if i + int(start) == pc {
            cfg.highlightPCColour.Printf("0x%04x", elem)
        } else {
            fmt.Printf("0x%04x", elem)
        }

        //Space between cells.
        fmt.Printf("  ")
    }

    fmt.Printf("\n")
}