package main

import (
    "fmt"
    "os"
    "regexp"
    "strings"
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
            fmt.Printf("idx %d check copy regex returned %d: %s\n", i, len(copy_match), strings.Join(copy_match, ","))
            env = append(env, fmt.Sprintf("%s=%s", copy_match[1], os.Getenv(copy_match[2])))
            continue
        }

        set_match := env_set_regex.FindStringSubmatch(os.Args[i])
        if len(set_match) > 0 {
            fmt.Printf("idx %d check set regex returned %d: %s\n", i, len(set_match), strings.Join(set_match, ","))
            env = append(env, fmt.Sprintf("%s=%s", set_match[1], set_match[2]))
            continue
        }

        fmt.Printf("Didn't match at elt %d\n", i)
        first_exec_idx = i
        break
    }

    fmt.Printf("Remaining arguments: %s\n", strings.Join(os.Args[first_exec_idx:], ","))
}
