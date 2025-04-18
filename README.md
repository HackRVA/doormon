# doormon

`doormon` listens to mqtt events and displays whether or not your access swipe was successful or not.

The current rfid reader at hackrva only provides positive feedback -- i.e. your swipe was successful and the door is unlocked.
It does this by powering an LED when the door lock gets power.  It does not provide negative feedback.  So, if your rfid fob isn't working or your subscription has expired, you might be at the door swiping multiple times (with no feedback) thinking that the reader isn't working.

`doormon` is a simple way to provide negative and positive feedback. 


> note: this is designed to run on some kind of potato that runs linux (like a raspberry pi zero 2 w). There are a few syscalls that are nix specific that won't work on windows or mac(maybe).

## Why no microcontroller?
A microcontroller would likely be a more reliable solution, but I basically have no idea what I'm doing with electronics most of the time. 
So, that would take much longer for me to develop.

I would be happy to contribute if someone else wants to head up a microcontroller-based solution.

## Testing

### Testing with mqtt
a docker compose file is included for running an mqtt broker.

```bash
# start the mqtt broker
make run-broker

# start doormon
doormon
## should display `waiting for activity`
```

> note: to run the mqtt scripts, you need [hivemq mqtt client](https://hivemq.github.io/mqtt-cli/docs/quick-start/)

in a separate terminal, you could run:
```bash
bash ./scripts/mqtt_pub.known.sh
bash ./scripts/mqtt_pub.unknown.sh
```

### Testing with fifo
If docker isn't your jam, the fifo mode is a more simple setup.

```bash
doormon --mode fifo
```

in another terminal, you could run
```bash
bash ./scripts/fifo.known.sh
bash ./scripts/fifo.unknown.sh
```

