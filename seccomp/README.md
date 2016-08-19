## seccomp

Seccomp provides application sandboxing mechanism in the Linux kernel.
Seccomp configuration allows one to configure actions to take for matched syscalls and furthermore also allows matching on values passed as arguments to syscalls.
For more information about Seccomp, see [Seccomp kernel documentation](https://www.kernel.org/doc/Documentation/prctl/seccomp_filter.txt)
The actions, architectures, and operators are strings that match the definitions in seccomp.h from [libseccomp](https://github.com/seccomp/libseccomp) and are translated to corresponding values.
A valid list of constants as of libseccomp v2.3.0 is shown below.

Architecture Constants
* `SCMP_ARCH_X86`
* `SCMP_ARCH_X86_64`
* `SCMP_ARCH_X32`
* `SCMP_ARCH_ARM`
* `SCMP_ARCH_AARCH64`
* `SCMP_ARCH_MIPS`
* `SCMP_ARCH_MIPS64`
* `SCMP_ARCH_MIPS64N32`
* `SCMP_ARCH_MIPSEL`
* `SCMP_ARCH_MIPSEL64`
* `SCMP_ARCH_MIPSEL64N32`
* `SCMP_ARCH_PPC`
* `SCMP_ARCH_PPC64`
* `SCMP_ARCH_PPC64LE`
* `SCMP_ARCH_S390`
* `SCMP_ARCH_S390X`

Action Constants:
* `SCMP_ACT_KILL`
* `SCMP_ACT_TRAP`
* `SCMP_ACT_ERRNO`
* `SCMP_ACT_TRACE`
* `SCMP_ACT_ALLOW`

Operator Constants:
* `SCMP_CMP_NE`
* `SCMP_CMP_LT`
* `SCMP_CMP_LE`
* `SCMP_CMP_EQ`
* `SCMP_CMP_GE`
* `SCMP_CMP_GT`
* `SCMP_CMP_MASKED_EQ`

###### Example

```json
"seccomp": {
    "defaultAction": "SCMP_ACT_ALLOW",
    "architectures": [
        "SCMP_ARCH_X86"
    ],
    "syscalls": [
        {
            "name": "getcwd",
            "action": "SCMP_ACT_ERRNO"
        }
    ]
}
```

### Significant syscalls blocked by the default profile

`sample.json` secccomp profile is a whitelist which specifies the calls that
are allowed. The table below lists the significant (but not all) syscalls that
are effectively blocked because they are not on the whitelist. The table includes
the reason each syscall is blocked rather than white-listed.

| Syscall             | Description                                                                                                                           |
|---------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `acct`              | Accounting syscall which could let containers disable their own resource limits or process accounting. Also gated by `CAP_SYS_PACCT`. |
| `add_key`           | Prevent containers from using the kernel keyring, which is not namespaced.                                   |
| `adjtimex`          | Similar to `clock_settime` and `settimeofday`, time/date is not namespaced.                                  |
| `bpf`               | Deny loading potentially persistent bpf programs into kernel, already gated by `CAP_SYS_ADMIN`.              |
| `clock_adjtime`     | Time/date is not namespaced.                                                                                 |
| `clock_settime`     | Time/date is not namespaced.                                                                                 |
| `clone`             | Deny cloning new namespaces. Also gated by `CAP_SYS_ADMIN` for CLONE_* flags, except `CLONE_USERNS`.         |
| `create_module`     | Deny manipulation and functions on kernel modules.                                                           |
| `delete_module`     | Deny manipulation and functions on kernel modules. Also gated by `CAP_SYS_MODULE`.                           |
| `finit_module`      | Deny manipulation and functions on kernel modules. Also gated by `CAP_SYS_MODULE`.                           |
| `get_kernel_syms`   | Deny retrieval of exported kernel and module symbols.                                                        |
| `get_mempolicy`     | Syscall that modifies kernel memory and NUMA settings. Already gated by `CAP_SYS_NICE`.                      |
| `init_module`       | Deny manipulation and functions on kernel modules. Also gated by `CAP_SYS_MODULE`.                           |
| `ioperm`            | Prevent containers from modifying kernel I/O privilege levels. Already gated by `CAP_SYS_RAWIO`.             |
| `iopl`              | Prevent containers from modifying kernel I/O privilege levels. Already gated by `CAP_SYS_RAWIO`.             |
| `kcmp`              | Restrict process inspection capabilities, already blocked by dropping `CAP_PTRACE`.                          |
| `kexec_file_load`   | Sister syscall of `kexec_load` that does the same thing, slightly different arguments.                       |
| `kexec_load`        | Deny loading a new kernel for later execution.                                                               |
| `keyctl`            | Prevent containers from using the kernel keyring, which is not namespaced.                                   |
| `lookup_dcookie`    | Tracing/profiling syscall, which could leak a lot of information on the host.                                |
| `mbind`             | Syscall that modifies kernel memory and NUMA settings. Already gated by `CAP_SYS_NICE`.                      |
| `mount`             | Deny mounting, already gated by `CAP_SYS_ADMIN`.                                                             |
| `move_pages`        | Syscall that modifies kernel memory and NUMA settings.                                                       |
| `name_to_handle_at` | Sister syscall to `open_by_handle_at`. Already gated by `CAP_SYS_NICE`.                                      |
| `nfsservctl`        | Deny interaction with the kernel nfs daemon.                                                                 |
| `open_by_handle_at` | Cause of an old container breakout. Also gated by `CAP_DAC_READ_SEARCH`.                                     |
| `perf_event_open`   | Tracing/profiling syscall, which could leak a lot of information on the host.                                |
| `personality`       | Prevent container from enabling BSD emulation. Not inherently dangerous, but poorly tested, potential for a lot of kernel vulns. |
| `pivot_root`        | Deny `pivot_root`, should be privileged operation.                                                           |
| `process_vm_readv`  | Restrict process inspection capabilities, already blocked by dropping `CAP_PTRACE`.                          |
| `process_vm_writev` | Restrict process inspection capabilities, already blocked by dropping `CAP_PTRACE`.                          |
| `ptrace`            | Tracing/profiling syscall, which could leak a lot of information on the host. Already blocked by dropping `CAP_PTRACE`. |
| `query_module`      | Deny manipulation and functions on kernel modules.                                                            |
| `quotactl`          | Quota syscall which could let containers disable their own resource limits or process accounting. Also gated by `CAP_SYS_ADMIN`. |
| `reboot`            | Don't let containers reboot the host. Also gated by `CAP_SYS_BOOT`.                                           |
| `request_key`       | Prevent containers from using the kernel keyring, which is not namespaced.                                    |
| `set_mempolicy`     | Syscall that modifies kernel memory and NUMA settings. Already gated by `CAP_SYS_NICE`.                       |
| `setns`             | Deny associating a thread with a namespace. Also gated by `CAP_SYS_ADMIN`.                                    |
| `settimeofday`      | Time/date is not namespaced. Also gated by `CAP_SYS_TIME`.                                                    |
| `stime`             | Time/date is not namespaced. Also gated by `CAP_SYS_TIME`.                                                    |
| `swapon`            | Deny start/stop swapping to file/device. Also gated by `CAP_SYS_ADMIN`.                                       |
| `swapoff`           | Deny start/stop swapping to file/device. Also gated by `CAP_SYS_ADMIN`.                                       |
| `sysfs`             | Obsolete syscall.                                                                                             |
| `_sysctl`           | Obsolete, replaced by /proc/sys.                                                                              |
| `umount`            | Should be a privileged operation. Also gated by `CAP_SYS_ADMIN`.                                              |
| `umount2`           | Should be a privileged operation.                                                                             |
| `unshare`           | Deny cloning new namespaces for processes. Also gated by `CAP_SYS_ADMIN`, with the exception of `unshare --user`. |
| `uselib`            | Older syscall related to shared libraries, unused for a long time.                                            |
| `userfaultfd`       | Userspace page fault handling, largely needed for process migration.                                          |
| `ustat`             | Obsolete syscall.                                                                                             |
| `vm86`              | In kernel x86 real mode virtual machine. Also gated by `CAP_SYS_ADMIN`.                                       |
| `vm86old`           | In kernel x86 real mode virtual machine. Also gated by `CAP_SYS_ADMIN`.                                       |