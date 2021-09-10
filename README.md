# socketcan


*socketcan* provides an interface to a [CAN bus](https://www.kernel.org/doc/Documentation/networking/can.txt) to read and write frames.
The library is based on the [SocketCAN](https://github.com/torvalds/linux/blob/master/include/uapi/linux/can.h) network stack on Linux.

## Features

* supports kernel timestamps for received frames
* receive and parse error frames
* `os.File` interface
* CAN filters
