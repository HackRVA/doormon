#!/bin/bash
FIFO="/tmp/door_notifier"

if [[ ! -p $FIFO ]]; then
    echo "FIFO $FIFO does not exist."
    exit 1
fi

echo "success" > "$FIFO"
echo "Sent known (success) message to $FIFO"
