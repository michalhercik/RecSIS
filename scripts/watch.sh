#!/bin/bash

(templ generate; go build -C ./mock_cas -o mockcas.exe; ./mock_cas/mockcas.exe --cert server.crt --key server.key) &
wgo -file=".go" -file=".templ" -xfile="_templ.go" templ generate :: go build -C ./webapp -o recsis :: ./webapp/recsis --config ./webapp/config.dev.toml