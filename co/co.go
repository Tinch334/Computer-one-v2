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
	case 6:
		retAddr = &(ci.regs.R6)
	case 7:
		retAddr = &(ci.regs.R7)
	default: //Avoid returning an error.
		retAddr = &(ci.regs.R0)
	}

	ci.setFlags(*retAddr)

	return retAddr
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
	ins := getInstruction(word)
	firstRegPtr := ci.getRegisterPtr(getFirstRegister(word))

	pcModified := false

	//The PC must increment by at least one every tick.
	ci.addPCinc()

	switch(ins) {
	//Load/store.
	case LD:
		b, regPtr, opr := ci.getRegisterOrImmediate(word)

		if b {
			*firstRegPtr = ci.GetMemoryCell(*regPtr)
		} else {
			*firstRegPtr = ci.GetMemoryCell(opr)
		}
	case ST:
		b, regPtr, opr := ci.getRegisterOrImmediate(word)

		if b {
			ci.SetMemoryCell(*regPtr, *firstRegPtr)
		} else {
			ci.SetMemoryCell(opr, *firstRegPtr)
		}

	case MOV:
		b, regPtr, opr := ci.getRegisterOrImmediate(word)

		if b {
			*firstRegPtr = *regPtr
		} else {
			*firstRegPtr = opr
		}
	
	//Arithmetic operations.
	case ADD:
		b, regPtr, opr := ci.getRegisterOrImmediate(word)

		if b {
			*firstRegPtr += *regPtr
		} else {
			*firstRegPtr += opr
		}
	
	case MUL:
		b, regPtr, opr := ci.getRegisterOrImmediate(word)

		if b {
			*firstRegPtr *= *regPtr
		} else {
			*firstRegPtr *= opr
		}
		
	//Logic operations.
	case AND:
		b, regPtr, opr := ci.getRegisterOrImmediate(word)

		if b {
			*firstRegPtr &= *regPtr
		} else {
			*firstRegPtr &= opr
		}

	case NOT:
		*firstRegPtr = ^(*firstRegPtr)
		
	case OR:
		b, regPtr, opr := ci.getRegisterOrImmediate(word)

		if b {
			*firstRegPtr |= *regPtr
		} else {
			*firstRegPtr |= opr
		}

	case SHL:
		b, regPtr, opr := ci.getRegisterOrImmediate(word)

		if b {
			*firstRegPtr *= *firstRegPtr << *regPtr
		} else {
			*firstRegPtr *= *firstRegPtr << opr
		}
	
	case SHR:
		b, regPtr, opr := ci.getRegisterOrImmediate(word)

		if b {
			*firstRegPtr *= *firstRegPtr >> *regPtr
		} else {
			*firstRegPtr *= *firstRegPtr >> opr
		}

    //Flow control.
	case JMP:
		b, _, operand := ci.getRegisterOrImmediate(word)

		//Invalid operand, do nothing.
		if b {
			return errors.New("Invalid operand for JMP"), true
		}

		f := ci.flags

		//Check jump flags.
		if (getBit(word, 2) && f.N) {
			ci.regs.PC = operand
			pcModified = true
		} else if (getBit(word, 1) && f.P) {
			ci.regs.PC = operand
			pcModified = true
		} else if (getBit(word, 0) && f.Z) {
			ci.regs.PC = operand
			pcModified = true
		}

	case JSR:
		b, _, operand := ci.getRegisterOrImmediate(word)

		//Invalid operand, do nothing.
		if b {
			return errors.New("Invalid operand for JSR"), true
		}

		ci.regs.R7 = (ci.regs.PC + 1) % MemorySize
		ci.regs.PC = operand
		pcModified = true

	case RET:
		ci.regs.PC = ci.regs.R7
		pcModified = true

	case NOP:
		
	case HLT:
		return nil, false
	}

	//Increment PC only if the instruction did not explicitly change it.
	if !pcModified {
		//Increment PC and check for overflow, reset increment counter.
		ci.regs.PC += ci.pcIncs
		ci.regs.PC = ci.regs.PC % MemorySize
		ci.pcIncs = 0
	}

	return nil, true
}