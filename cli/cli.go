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
    }

    run(reader, &control, &config)
}

func run(reader *bufio.Reader, ctrl *interpreterControl, cfg *interpreterConfig) {
    processInput(reader, ctrl, cfg)
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

    case HELP:
        fallthrough
    case HELP_SHORT:
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