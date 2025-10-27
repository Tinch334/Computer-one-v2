package co


import (
	"errors"
)


const MemorySize = 1024 //In words.

//The registers and flags are a separate structure to be able to return them.
type Registers struct {
	PC, R0, R1, R2, R3, R4, R5, R6, R7, RR uint16
}

type Flags struct {
	N, P, Z bool
}

//We use 16 bit words.
type ComputerInfo struct {
	//Registers.
	regs Registers

	//Status flags.
	flags Flags

	//How many times the PC must be incremented in the next tick.
	pcIncs uint16

	//Array representing memory.
	memory [MemorySize]uint16
}

const (
	LD = iota
	ST
	MOV
	ADD
	MUL
	AND
	NOT
	OR
	SHL
	SHR
	JMP
	JSR
	RET
	NOP
	HLT
)


func NewComputerInfo() *ComputerInfo {
	ci := ComputerInfo{
		regs: Registers{},
		flags: Flags{},
		pcIncs: 0,
		memory: [MemorySize]uint16{},
	}

	return &ci
}

/*
	OUTPUT FUNCTIONS
*/
func (ci *ComputerInfo) GetRegisters() Registers {
	return ci.regs
}

func (ci *ComputerInfo) GetFlags() Flags {
	return ci.flags
}

func (ci *ComputerInfo) GetMemory(start uint16, end uint16) (error, []uint16) {
	if start >= end {
		return errors.New("Invalid memory slice: start must be < end"), nil
	}
	if int(end) > MemorySize {
		return errors.New("Invalid memory slice: end out of bounds"), nil
	}

	returnedMemory := make([]uint16, end - start)

	for i := start; i < end; i++ {
		returnedMemory[i] = ci.memory[i + start]
	}

	return nil, returnedMemory
}

/*
	SETTING FUNCTIONS
*/
func (ci *ComputerInfo) setFlags(res uint16) {
	s := int16(res)

	ci.flags.N = s < 0
	ci.flags.P = s > 0
	ci.flags.Z = s == 0
}


//Sets the memory cells in the specified interval.
func (ci *ComputerInfo) SetMemory(start int, mem []uint16) error {
	if start + len(mem) > MemorySize {
		return errors.New("Invalid start position and memory length")
	}

	for i, elem := range mem {
		ci.memory[i + start] = elem
	}

	return nil
}

//Sets all CPU registers.
func (ci *ComputerInfo) SetRegisters(regs Registers, flags Flags) {
	ci.regs = regs
	ci.flags = flags
}

/*
	REGISTER INSTRUCTIONS
*/
//Takes an instruction, if it's in immediate mode returns the value and "false", otherwise "true" and a pointer to the appropriate register.
func (ci *ComputerInfo) getRegisterOrImmediate(ins uint16) (bool, *uint16, uint16) {
	//Check if double mode is enabled, if so load data from next memory cell.
	if getLowerByte(ins) == 0xFF {
		ci.addPCinc()
		return false, nil, ci.memory[ci.regs.PC + 1]
	}

	//Check immediate flag.
	if getBit(ins, 7) {
		regNum := getSecondRegister(ins)
		reg := ci.getRegisterPtr(regNum)

		return true, reg, 0
	}

	imm := getImmediate(ins)

	return false, nil, imm
}

//Adds one increment to the PC in the next tick.
func (ci *ComputerInfo) addPCinc() {
	ci.pcIncs += 1
}

/*
	INSTRUCTION INFORMATION
*/
func getInstruction(ins uint16) uint16 {
	return (ins & 0xF800) >> 12
}

func getFirstRegister(ins uint16) uint16 {
	return (ins & 0x0700) >> 8
}

func getSecondRegister(ins uint16) uint16 {
	return ins & 0x0003
}

func getImmediate(ins uint16) uint16 {
	return ins & 0x007F
}

func getLowerByte(ins uint16) uint16 {
	return ins & 0x00FF
}