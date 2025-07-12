package ui

import (
	"syscall"
	"unsafe"
)

type termios struct {
	Iflag  uint64
	Oflag  uint64
	Cflag  uint64
	Lflag  uint64
	Cc     [20]uint8
	Ispeed uint64
	Ospeed uint64
}

func makeRaw() (*termios, error) {
	var oldState termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, 0, syscall.TIOCGETA, uintptr(unsafe.Pointer(&oldState)))
	if errno != 0 {
		return nil, errno
	}

	newState := oldState
	newState.Lflag &^= syscall.ECHO | syscall.ICANON
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, 0, syscall.TIOCSETA, uintptr(unsafe.Pointer(&newState)))
	if errno != 0 {
		return nil, errno
	}

	return &oldState, nil
}

func restore(oldState *termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, 0, syscall.TIOCSETA, uintptr(unsafe.Pointer(oldState)))
	if errno != 0 {
		return errno
	}
	return nil
}
