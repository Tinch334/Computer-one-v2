package co

import (
	"errors"
)

/*
	INTERNAL INSTRUCTION PROCESSING
*/
//Returns the n-th bit from a 16 bit unsigned integer, 0 corresponds to the lsb.
func getBit(ins uint16, pos int) bool {
	return (ins & (1 << pos)) != 0
}

//Returns a pointer to the appropriate general register, if the argument is invalid the first register is returned.
func (ci *ComputerInfo) getRegisterPtr(arg uint16) *uint16 {
	var retAddr *uint16

	switch(arg) {
	case 0:
		retAddr = &(ci.regs.R0)
	case 1:
		retAddr = &(ci.regs.R1)
	case 2:
		retAddr = &(ci.regs.R2)
	case 3:
		retAddr = &(ci.regs.R3)
	case 4:
		retAddr = &(ci.regs.R4)
	case 5:
		retAddr = &(ci.regs.R5)
	default: //Avoid returning an error.
		retAddr = &(ci.regs.R0)
	}

	ci.setFlags(*retAddr)

	return retAddr
}

/*
	INSTRUCTION INFORMATION
*/
func GetInstruction(ins uint16) uint16 {
	return (ins & 0xF000) >> 12
}

func GetFirstRegister(ins uint16) uint16 {
	return (ins & 0x0700) >> 8
}

func GetSecondRegister(ins uint16) uint16 {
	return ins & 0x0003
}

func GetImmediate(ins uint16) uint16 {
	return ins & 0x00FF
}

//Takes an instruction, if it's in immediate mode returns the value and "false", otherwise "true" and a pointer to the appropriate register.
func getRegisterOrImmediate(ci *ComputerInfo, ins uint16) (bool, *uint16, uint16) {
	//Check immediate flag.
	if getBit(ins, 11) {
		regNum := GetSecondRegister(ins)
		reg := ci.getRegisterPtr(regNum)

		return true, reg, 0
	}

	imm := GetImmediate(ins)
	return false, nil, imm
}

//Performs left shift the specified amount.
func leftShift(value uint16, amount uint16) uint16 {
	return value << amount
}

//Performs logic right shift the specified amount.
func rightShift(value uint16, amount uint16) uint16 {
	return value >> amount
}


/*
	INTERPRETER
*/
func (ci *ComputerInfo) Step() (error, bool) {
	word := ci.memory[int(ci.regs.PC)]
	ins := GetInstruction(word)
	firstRegPtr := ci.getRegisterPtr(GetFirstRegister(word))

	ci.regs.PC++
	//Check for program counter exceeding memory size.
	if ci.regs.PC >= MemorySize {
		ci.regs.PC = 0
	}

	switch(ins) {
	//Load/store.
	case LD:
		b, regPtr, opr := getRegisterOrImmediate(ci, word)

		if b {
			if *regPtr >= MemorySize {
				return errors.New("Invalid value for LD"), true
			}

			*firstRegPtr = ci.memory[*regPtr]
		} else {
			*firstRegPtr = ci.memory[opr]
		}
	
	case ST:
		b, regPtr, opr := getRegisterOrImmediate(ci, word)

		if b {
			if *regPtr >= MemorySize {
				return errors.New("Invalid value for ST"), true
			}

			ci.memory[*regPtr] = *firstRegPtr
		} else {
			ci.memory[opr] = *firstRegPtr
		}

	case MOV:
		b, regPtr, opr := getRegisterOrImmediate(ci, word)

		if b {
			*firstRegPtr = *regPtr
		} else {
			*firstRegPtr = opr
		}
	
	//Arithmetic operations.
	case ADD:
		b, regPtr, opr := getRegisterOrImmediate(ci, word)

		if b {
			*firstRegPtr += *regPtr
		} else {
			*firstRegPtr += opr
		}
	
	case MUL:
		b, regPtr, opr := getRegisterOrImmediate(ci, word)

		if b {
			*firstRegPtr *= *regPtr
		} else {
			*firstRegPtr *= opr
		}
		
	//Logic operations.
	case AND:
		b, regPtr, opr := getRegisterOrImmediate(ci, word)

		if b {
			*firstRegPtr &= *regPtr
		} else {
			*firstRegPtr &= opr
		}

	case NOT:
		*firstRegPtr = ^(*firstRegPtr)
		
	case OR:
		b, regPtr, opr := getRegisterOrImmediate(ci, word)

		if b {
			*firstRegPtr |= *regPtr
		} else {
			*firstRegPtr |= opr
		}

	case SHL:
		b, regPtr, opr := getRegisterOrImmediate(ci, word)

		if b {
			*firstRegPtr *= *firstRegPtr << *regPtr
		} else {
			*firstRegPtr *= *firstRegPtr << opr
		}
	
	case SHR:
		b, regPtr, opr := getRegisterOrImmediate(ci, word)

		if b {
			*firstRegPtr *= *firstRegPtr >> *regPtr
		} else {
			*firstRegPtr *= *firstRegPtr >> opr
		}

    //Flow control.
	case JMP:
		b, _, operand := getRegisterOrImmediate(ci, word)

		//Invalid operand, do nothing.
		if b {
			return errors.New("Invalid operand for JMP"), true
		}

		f := ci.flags

		//Check jump flags.
		if (getBit(word, 10) && f.N) {
			ci.regs.PC = operand
		} else if (getBit(word, 9) && f.P) {
			ci.regs.PC = operand
		} else if (getBit(word, 8) && f.Z) {
			ci.regs.PC = operand
		}

	case JSR:
		b, _, operand := getRegisterOrImmediate(ci, word)

		//Invalid operand, do nothing.
		if b {
			return errors.New("Invalid operand for JSR"), true
		}

		ci.regs.RR = ci.regs.PC
		ci.regs.PC = operand

	case RET:
		ci.regs.PC = ci.regs.RR

	case NOP:
		
	case HLT:
		return nil, false
	}

	return nil, true
}