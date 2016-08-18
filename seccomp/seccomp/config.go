package seccomp

import "fmt"

var operators = map[string]Operator{
	"SCMP_CMP_NE":        NotEqualTo,
	"SCMP_CMP_LT":        LessThan,
	"SCMP_CMP_LE":        LessThanOrEqualTo,
	"SCMP_CMP_EQ":        EqualTo,
	"SCMP_CMP_GE":        GreaterThanOrEqualTo,
	"SCMP_CMP_GT":        GreaterThan,
	"SCMP_CMP_MASKED_EQ": MaskEqualTo,
}

var actions = map[string]Action{
	"SCMP_ACT_KILL":  Kill,
	"SCMP_ACT_ERRNO": Errno,
	"SCMP_ACT_TRAP":  Trap,
	"SCMP_ACT_ALLOW": Allow,
	"SCMP_ACT_TRACE": Trace,
}

var archs = map[string]string{
	"SCMP_ARCH_X86":         "x86",
	"SCMP_ARCH_X86_64":      "amd64",
	"SCMP_ARCH_X32":         "x32",
	"SCMP_ARCH_ARM":         "arm",
	"SCMP_ARCH_AARCH64":     "arm64",
	"SCMP_ARCH_MIPS":        "mips",
	"SCMP_ARCH_MIPS64":      "mips64",
	"SCMP_ARCH_MIPS64N32":   "mips64n32",
	"SCMP_ARCH_MIPSEL":      "mipsel",
	"SCMP_ARCH_MIPSEL64":    "mipsel64",
	"SCMP_ARCH_MIPSEL64N32": "mipsel64n32",
}

// ConvertStringToOperator converts a string into a Seccomp comparison operator.
// Comparison operators use the names they are assigned by Libseccomp's header.
// Attempting to convert a string that is not a valid operator results in an
// error.
func ConvertStringToOperator(in string) (Operator, error) {
	if op, ok := operators[in]; ok == true {
		return op, nil
	}
	return 0, fmt.Errorf("string %s is not a valid operator for seccomp", in)
}

// ConvertStringToAction converts a string into a Seccomp rule match action.
// Actions use the names they are assigned in Libseccomp's header, though some
// (notable, SCMP_ACT_TRACE) are not available in this implementation and will
// return errors.
// Attempting to convert a string that is not a valid action results in an
// error.
func ConvertStringToAction(in string) (Action, error) {
	if act, ok := actions[in]; ok == true {
		return act, nil
	}
	return 0, fmt.Errorf("string %s is not a valid action for seccomp", in)
}

// ConvertStringToArch converts a string into a Seccomp comparison arch.
func ConvertStringToArch(in string) (string, error) {
	if arch, ok := archs[in]; ok == true {
		return arch, nil
	}
	return "", fmt.Errorf("string %s is not a valid arch for seccomp", in)
}
