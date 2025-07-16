#!/bin/bash

wgo \
    -file=".go" -file=".templ" \
    -xfile="_templ.go" \
    templ generate :: \
    go build -C ./webapp -o recsis :: \
    ./webapp/recsis --config ./webapp/config.dev.toml