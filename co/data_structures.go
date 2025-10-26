package co


import (
	"errors"
)


const MemorySize = 1024 //In words.

//The registers and flags are a separate structure to be able to return them.
type Registers struct {
	PC, R0, R1, R2, R3, R4, R5, RR uint16
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

func (ci *ComputerInfo) SetMemory(start int, mem []uint16) error {
	if start + len(mem) > MemorySize {
		return errors.New("Invalid start position and memory length")
	}

	for i, elem := range mem {
		ci.memory[i + start] = elem
	}

	return nil
}