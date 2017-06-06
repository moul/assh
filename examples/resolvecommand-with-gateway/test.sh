#!/bin/sh

assh -D -c ./assh.yml connect --dry-run a/b/c 2>&1 | cut -d\  -f 2- > test.log || true
