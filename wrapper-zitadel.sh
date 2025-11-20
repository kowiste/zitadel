#!/bin/bash
set -e

echo "--- 1. Starting ZITADEL Server ---"
# Start ZITADEL in background
# Note: We call /app/zitadel directly as we copied it there
/app/zitadel start-from-init --masterkey "${ZITADEL_MASTERKEY}" --tlsMode disabled --config /zitadel-init.yaml &
ZITADEL_PID=$!

echo "--- 2. Waiting for ZITADEL to be ready ---"
TIMEOUT=300
COUNTER=0
# ZITADEL health endpoint
HEALTH_URL="http://localhost:8080/debug/healthz"

while ! curl -f -s ${HEALTH_URL} > /dev/null 2>&1; do
    if [ $COUNTER -ge $TIMEOUT ]; then
        echo "Timeout waiting for ZITADEL."
        exit 1
    fi
    echo "Waiting for ZITADEL... ($COUNTER/$TIMEOUT)"
    sleep 5
    COUNTER=$((COUNTER + 5))
done

echo "--- 3. Running Setup Script ---"
# Set environment variables for the setup script
export ZITADEL_PAT_KEY_PATH="/zitadel-pat-key/admin.pat"

# Execute the setup script
if ! /usr/local/bin/setup_zitadel.sh; then
    echo "Setup script failed! Killing ZITADEL."
    kill $ZITADEL_PID
    exit 1
fi

echo "--- 4. Setup Complete. Running indefinitely ---"
wait $ZITADEL_PID