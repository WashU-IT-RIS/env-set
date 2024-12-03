package main

import (
    "fmt"
    "os"
    "regexp"
    "strings"
    "os/exec"
    "syscall"
    "errors"
)

func main() {
    env_copy_regex := regexp.MustCompile(`(\w+)=\$(\w+)`)
    env_set_regex  := regexp.MustCompile(`(\w+)=(\w+)`)

    // The first argv index that's the executable we were asked to run
    // in the container.  Subsequent args are params to this executable.
    first_exec_idx := -1

    env := os.Environ()

    for i := 1; i < len(os.Args); i++ {
        copy_match := env_copy_regex.FindStringSubmatch(os.Args[i])
        if len(copy_match) > 0 {
            //fmt.Printf("idx %d check copy regex returned %d: %s\n", i, len(copy_match), strings.Join(copy_match, ","))
            env = append(env, fmt.Sprintf("%s=%s", copy_match[1], os.Getenv(copy_match[2])))
            continue
        }

        set_match := env_set_regex.FindStringSubmatch(os.Args[i])
        if len(set_match) > 0 {
            //fmt.Printf("idx %d check set regex returned %d: %s\n", i, len(set_match), strings.Join(set_match, ","))
            env = append(env, fmt.Sprintf("%s=%s", set_match[1], set_match[2]))
            continue
        }

        //fmt.Printf("Didn't match at elt %d\n", i)
        first_exec_idx = i
        break
    }

    //fmt.Printf("Remaining arguments: %s\n", strings.Join(os.Args[first_exec_idx:], ","))

    // If the executable pathname contains a "/", then it's meant to be
    // an explicit path.  If no "/", then find it in the PATH.
    slash_idx := strings.Index(os.Args[first_exec_idx], "/")
    var bin_path string
    if slash_idx == -1 {
        var look_err error
        bin_path, look_err = exec.LookPath(os.Args[first_exec_idx])
        if look_err != nil {
            fmt.Fprintf(os.Stderr, "%s: command not found: %v\n", os.Args[first_exec_idx], look_err)
            os.Exit(127) // Matches exit code for bash for "Command not found"
        }
    } else {
        bin_path = os.Args[first_exec_idx]
    }

    // execvp() the requested process
    exec_err := syscall.Exec(bin_path, os.Args[first_exec_idx:], env)
    if exec_err != nil {
        fmt.Printf("Can't execute %s: %v\n", os.Args[first_exec_idx], exec_err)
        if _, stat_err := os.Stat(bin_path); stat_err == nil {
            os.Exit(126)  // Same as bash, path exists but is not executable
        } else if errors.Is(stat_err, os.ErrNotExist) {
            os.Exit(127)
        } else {
            // Not sure what the conditions would be to get here
            os.Exit(126)
        }
    }
    // Shouldn't get here.  Either the exec replaces this process, or there
    // was an error caught by the above if/else
}
