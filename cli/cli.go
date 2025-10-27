package cli

import (
    "fmt"
    "strings"
    "bufio"
    "os"

    "github.com/Tinch334/Computer-one-v2/co"
    "github.com/fatih/color"
)

func RunCli() {
    ci := co.NewComputerInfo()

    memLoad := []uint16{
        0b0000010000000011, //LD r4 <- mem[3]
        0b0000001010000100, //LD r2 <- mem[r2]
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
    }

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

    if len(line) == 0 {
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
        breakpointHandler(ctrl, cfg, arguments)

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
        configurationHandler(cfg, arguments)

    default:
        fmt.Println("Unknown command, use \"h\" for help")
        return
    }
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

    fmt.Printf("PC: 0x%04x | NPZ: %s | RR: 0x%04x\n",
        regs.PC, flagsStr, regs.RR)

    fmt.Printf("R0: 0x%04x R1: 0x%04x R2: 0x%04x R3:0x%04x R4:0x%04x R5:0x%04x R6:0x%04x R7:0x%04x\n",
        regs.R0, regs.R1, regs.R2, regs.R3, regs.R4, regs.R5, regs.R6, regs.R7)
}

//Prints memory contents, note that "tabwriter" cannot be used because ANSI escape codes are used for colour, and they get counted
//by the package.
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