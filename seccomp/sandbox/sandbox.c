#define _GNU_SOURCE 1
#include <stdio.h>
#include <stddef.h>
#include <stdlib.h>
#include <unistd.h>

#include "seccomp-bpf.h"
// Added for debugging.
#include "syscall-reporter.h"

static int install_syscall_filter(void)
{
        struct sock_filter filter[] = {
                /* Validate architecture. */
                VALIDATE_ARCHITECTURE,
                /* Grab the system call number. */
                EXAMINE_SYSCALL,
                /* List allowed syscalls. */
                ALLOW_SYSCALL(rt_sigreturn),
                #ifdef __NR_sigreturn
                ALLOW_SYSCALL(sigreturn),
                #endif
                ALLOW_SYSCALL(exit_group),
                ALLOW_SYSCALL(exit),
                ALLOW_SYSCALL(read),
                ALLOW_SYSCALL(write),
                /// New
                ALLOW_SYSCALL(fstat),
                ALLOW_SYSCALL(mmap),
                ALLOW_SYSCALL(brk),
                ALLOW_SYSCALL(nanosleep),
                ALLOW_SYSCALL(execve),
                ALLOW_SYSCALL(access),
                ALLOW_SYSCALL(open),
                ALLOW_SYSCALL(close),
                ALLOW_SYSCALL(mprotect),
                ALLOW_SYSCALL(arch_prctl),
                ALLOW_SYSCALL(munmap),
                // Enable only for minimal bash scripts.
                ALLOW_SYSCALL(getuid),
                ALLOW_SYSCALL(getgid),
                ALLOW_SYSCALL(geteuid),
                ALLOW_SYSCALL(getegid),
                ALLOW_SYSCALL(rt_sigprocmask),
                ALLOW_SYSCALL(sysinfo),
                ALLOW_SYSCALL(rt_sigaction),
                ALLOW_SYSCALL(uname),
                ALLOW_SYSCALL(stat),
                ALLOW_SYSCALL(getpid),
                ALLOW_SYSCALL(getppid),
                ALLOW_SYSCALL(getpgrp),
                ALLOW_SYSCALL(getrlimit),
                ALLOW_SYSCALL(sched_getattr),
                ALLOW_SYSCALL(ioctl),
                ALLOW_SYSCALL(lseek),
                ALLOW_SYSCALL(fcntl),
                ALLOW_SYSCALL(dup2),
                // Allowed for nodejs.
                ALLOW_SYSCALL(set_tid_address),
                ALLOW_SYSCALL(set_robust_list),
                ALLOW_SYSCALL(setrlimit),
                ALLOW_SYSCALL(pipe2),
                ALLOW_SYSCALL(clock_getres),
                ALLOW_SYSCALL(epoll_create1),
                ALLOW_SYSCALL(eventfd2),
                ALLOW_SYSCALL(prctl),
                ALLOW_SYSCALL(poll),
                ALLOW_SYSCALL(dup3),
                ALLOW_SYSCALL(madvise),
                // Allowed for python.
                ALLOW_SYSCALL(readlink),
                ALLOW_SYSCALL(lstat),
                ALLOW_SYSCALL(getdents),
                ALLOW_SYSCALL(getcwd),
                ALLOW_SYSCALL(futex),
                ALLOW_SYSCALL(select),
                // Allowed for Go.
                ALLOW_SYSCALL(sched_getaffinity),
                ALLOW_SYSCALL(sigaltstack),
                ALLOW_SYSCALL(gettid),
                ALLOW_SYSCALL(clone),
                KILL_PROCESS,                
        };
        struct sock_fprog prog = {
                .len = (unsigned short)(sizeof(filter)/sizeof(filter[0])),
                .filter = filter,
        };
        
        if (prctl(PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0)) {
                perror("prctl(NO_NEW_PRIVS)");
                goto failed;
        }
        if (prctl(PR_SET_SECCOMP, SECCOMP_MODE_FILTER, &prog)) {
                perror("prctl(SECCOMP)");
                goto failed;
        }
        return 0;
        
failed:
        if (errno == EINVAL)
                fprintf(stderr, "SECCOMP_FILTER is not available. :(\n");
        return 1;
}

int main(int argc, char *argv[])
{
        if (install_syscall_filter()) 
                return 1;

        char* new_argv[1024];
        int i = execv(argv[1], new_argv);
        if (i > 0) {
                printf("Do not know what to expect %d", i);
                return 1;
        }
        // Success.
	return 0;
}
