# state-logger

[![Build Status](https://travis-ci.org/posener/state-logger.svg?branch=master)](https://travis-ci.org/posener/state-logger)

A package used to log state changes.

Supports:
* Log using your favorite log library
* Set log interval
* Different log levels for errors and for success

## Usage

Sometimes you may have a function that fails.
The code below shows a loop that tries to connect

```go
var c Connection
for {
    c, err = connect()
    if err == nil {
        break
    }
    l.Println(err)
    time.Sleep(10 * time.Second)
}
```

If the connection fails, this function returns an error, and
your logs will get an error message every 10 seconds.
This makes them ugly, and also harder to debug.

This is the place you would like to use the state-logger:

```go
var c Connection
sl := NewStateLogger("connection", log.Printf)
for {
    c, err = connect()
    if err == nil {
        break
    }
    sl.LogError(err)
    time.Sleep(10 * time.Second)
}
sl.Fixed()
```

This could make your logs from:

```
Failed connecting to server...
Failed connecting to server...
Failed connecting to server...
Failed connecting to server...
Failed connecting to server...
```

To:
```
[connection] error: Failed connecting to server... 
[connection] fixed!
```

In case that several errors occurred:
```
[connection] error: Failed connecting to server... 
[connection] error: DNS lookup failed
[connection] fixed!
```
