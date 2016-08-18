// +build linux,cgo

package seccomp

import (
	"fmt"
	"syscall"

	libseccomp "github.com/seccomp/libseccomp-golang"
)

var (
	actAllow = libseccomp.ActAllow
	actTrap  = libseccomp.ActTrap
	actKill  = libseccomp.ActKill
	actTrace = libseccomp.ActTrace.SetReturnCode(int16(syscall.EPERM))
	actErrno = libseccomp.ActErrno.SetReturnCode(int16(syscall.EPERM))

	// SeccompModeFilter refers to the syscall argument SECCOMP_MODE_FILTER.
	SeccompModeFilter = uintptr(2)
)

// Filters given syscalls in a container, preventing them from being used
// Started in the container init process, and carried over to all child processes
// Setns calls, however, require a separate invocation, as they are not children
// of the init until they join the namespace
func InitSeccomp(config *Seccomp) error {
	if config == nil {
		return fmt.Errorf("cannot initialize Seccomp - nil config passed")
	}

	defaultAction, err := getAction(config.DefaultAction)
	if err != nil {
		fmt.Println(config.DefaultAction)
		return fmt.Errorf("error initializing seccomp - invalid default action")
	}

	filter, err := libseccomp.NewFilter(defaultAction)
	if err != nil {
		return fmt.Errorf("error creating filter: %s", err)
	}

	// Add extra architectures
	for _, arch := range config.Architectures {
		scmpArch, err := libseccomp.GetArchFromString(arch)
		if err != nil {
			return err
		}

		if err := filter.AddArch(scmpArch); err != nil {
			return err
		}
	}

	// Unset no new privs bit
	if err := filter.SetNoNewPrivsBit(false); err != nil {
		return fmt.Errorf("error setting no new privileges: %s", err)
	}

	// Add a rule for each syscall
	for _, call := range config.Syscalls {
		if call == nil {
			return fmt.Errorf("encountered nil syscall while initializing Seccomp")
		}

		if err = matchCall(filter, call); err != nil {
			return err
		}
	}

	if err = filter.Load(); err != nil {
		return fmt.Errorf("error loading seccomp filter into kernel: %s", err)
	}

	return nil
}

// IsEnabled returns if the kernel has been configured to support seccomp.
func IsEnabled() bool {
	// Check if Seccomp is supported, via CONFIG_SECCOMP.
	_, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_GET_SECCOMP, 0, 0)
	if err != syscall.EINVAL {
		// Make sure the kernel has CONFIG_SECCOMP_FILTER.
		_, _, err = syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_SET_SECCOMP, SeccompModeFilter, 0)
		if err != syscall.EINVAL {
			return true
		}
	}
	return false
}

// Convert Libcontainer Action to Libseccomp ScmpAction
func getAction(act Action) (libseccomp.ScmpAction, error) {
	switch act {
	case Kill:
		return actKill, nil
	case Errno:
		return actErrno, nil
	case Trap:
		return actTrap, nil
	case Allow:
		return actAllow, nil
	case Trace:
		return actTrace, nil
	default:
		return libseccomp.ActInvalid, fmt.Errorf("invalid action, cannot use in rule")
	}
}

// Convert Libcontainer Operator to Libseccomp ScmpCompareOp
func getOperator(op Operator) (libseccomp.ScmpCompareOp, error) {
	switch op {
	case EqualTo:
		return libseccomp.CompareEqual, nil
	case NotEqualTo:
		return libseccomp.CompareNotEqual, nil
	case GreaterThan:
		return libseccomp.CompareGreater, nil
	case GreaterThanOrEqualTo:
		return libseccomp.CompareGreaterEqual, nil
	case LessThan:
		return libseccomp.CompareLess, nil
	case LessThanOrEqualTo:
		return libseccomp.CompareLessOrEqual, nil
	case MaskEqualTo:
		return libseccomp.CompareMaskedEqual, nil
	default:
		return libseccomp.CompareInvalid, fmt.Errorf("invalid operator, cannot use in rule")
	}
}

// Convert Libcontainer Arg to Libseccomp ScmpCondition
func getCondition(arg *Arg) (libseccomp.ScmpCondition, error) {
	cond := libseccomp.ScmpCondition{}

	if arg == nil {
		return cond, fmt.Errorf("cannot convert nil to syscall condition")
	}

	op, err := getOperator(arg.Op)
	if err != nil {
		return cond, err
	}

	return libseccomp.MakeCondition(arg.Index, op, arg.Value, arg.ValueTwo)
}

// Add a rule to match a single syscall
func matchCall(filter *libseccomp.ScmpFilter, call *Syscall) error {
	if call == nil || filter == nil {
		return fmt.Errorf("cannot use nil as syscall to block")
	}

	if len(call.Name) == 0 {
		return fmt.Errorf("empty string is not a valid syscall")
	}

	// If we can't resolve the syscall, assume it's not supported on this kernel
	// Ignore it, don't error out
	callNum, err := libseccomp.GetSyscallFromName(call.Name)
	if err != nil {
		return nil
	}

	// Convert the call's action to the libseccomp equivalent
	callAct, err := getAction(call.Action)
	if err != nil {
		return err
	}

	// Unconditional match - just add the rule
	if len(call.Args) == 0 {
		if err = filter.AddRule(callNum, callAct); err != nil {
			return err
		}
	} else {
		// Conditional match - convert the per-arg rules into library format
		conditions := []libseccomp.ScmpCondition{}

		for _, cond := range call.Args {
			newCond, err := getCondition(cond)
			if err != nil {
				return err
			}

			conditions = append(conditions, newCond)
		}

		if err = filter.AddRuleConditional(callNum, callAct, conditions); err != nil {
			return err
		}
	}

	return nil
}
