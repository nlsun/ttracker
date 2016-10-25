# the Trivial Tracker

Tracks CPU usage of processes over time.

External systems can send some JSON to the `/hook` endpoint. This then gets
logged along with a timestamp. This can be used to correlate events with the
periodic polling.
