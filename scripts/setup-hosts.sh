#!/bin/bash

HOSTS_LINE="127.0.0.1 mockcas"
HOSTS_FILE="/etc/hosts"

if ! grep -q "$HOSTS_LINE" "$HOSTS_FILE"; then
    echo "$HOSTS_LINE" | sudo tee -a "$HOSTS_FILE" > /dev/null
fi