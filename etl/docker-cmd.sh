#!/bin/sh

chmod 600 /root/.ssh/id_rsa

autossh -f -M 0 -gNC \
-o "ExitOnForwardFailure=yes" \
-o "ServerAliveInterval=10" \
-o "ServerAliveCountMax=3" \
-L 10502:mffout.karlov.mff.cuni.cz:10502 $ACHERON_USER@acheron.ms.mff.cuni.cz \
-p 42049 \
-o StrictHostKeyChecking=no \
-o UserKnownHostsFile=/dev/null

timeout=30
while ! nc -z localhost 10502; do
    sleep 1
    timeout=$((timeout - 1))
    if [ $timeout -le 0 ]; then
        echo "SSH tunnel setup timed out."
        exit 1
    fi
done

/etl --config /app/config.docker.toml

pkill autossh
