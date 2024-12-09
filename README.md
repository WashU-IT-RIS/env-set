envset
======

A clone for the POSIX utility `env` written in Go in order to be statically
compiled.

Description
-----------

LSF jobs that utilize GPUs need certain environment variables set to control
their execution.  In particular, `NVIDIA_VISIBLE_DEVICES` needs to be set to
the same value as `CUDA_VISIBLE_DEVICES`, but for some reason setting it using
normal channels like the `-e` or `--env-file` options to `docker run` do not
work.  This utility allows setting it in a transparent manner before running
whatever process the job actually intended to run, and should be safe to map
into the container and run since it is statically linked.

Usage
-----

envset [NAME=VALUE]... [NAME=$OTHER]... COMMAND [ARGS]...

Run the given COMMAND with its arguments with an altered environment by
setting the NAMEd variables to the given VALUEs, or by copying the value
of OTHER existing variables to the given NAMEs.
