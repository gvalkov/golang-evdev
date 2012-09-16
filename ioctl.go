package evdev

import (
	"syscall"
	"unsafe"
)

const (
	_IOC_NONE  = 0x0
	_IOC_WRITE = 0x1
	_IOC_READ  = 0x2

	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_DIRBITS  = 2
	_IOC_NRSHIFT  = 0

	_IOC_TYPESHIFT = _IOC_NRSHIFT   + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS
)

func _IOC(dir int, t int, nr int, size int) int {
	return  (dir << _IOC_DIRSHIFT) |    (t << _IOC_TYPESHIFT) |
		     (nr << _IOC_NRSHIFT)  | (size << _IOC_SIZESHIFT)
}

func _IOR(t int, nr int, size int) int {
	return _IOC(_IOC_READ, t, nr, size)
}

func ioctl(fd uintptr, name int, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, uintptr(name), uintptr(data))
	return err
}


// input.h ioctls
var EVIOCGID   = _IOR('E', 0x02, 8)  // 8 <- sizeof(struct input_id)
var EVIOCGNAME = _IOC(_IOC_READ, 'E', 0x06, MAX_NAME_SIZE)
var EVIOCGPHYS = _IOC(_IOC_READ, 'E', 0x07, MAX_NAME_SIZE)

func EVIOCGBIT(ev, len int) int {
	return _IOC(_IOC_READ, 'E', 0x20 + ev, len)
}
