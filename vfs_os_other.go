//go:build !windows && !linux && !darwin

package sqlite3

import (
	"io"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

func (vfsOSMethods) Sync(file *os.File, fullsync, dataonly bool) error {
	return file.Sync()
}

func (vfsOSMethods) Allocate(file *os.File, size int64) error {
	off, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	if size <= off {
		return nil
	}
	return file.Truncate(size)
}

func (vfsOSMethods) unlock(file *os.File, start, len int64) xErrorCode {
	if start == 0 && len == 0 {
		err := unix.Flock(int(file.Fd()), unix.LOCK_UN)
		if err != nil {
			return IOERR_UNLOCK
		}
	}
	return _OK
}

func (vfsOSMethods) readLock(file *os.File, start, len int64, timeout time.Duration) xErrorCode {
	var err error
	for {
		err = unix.Flock(int(file.Fd()), unix.LOCK_SH|unix.LOCK_NB)
		if errno, _ := err.(unix.Errno); errno != unix.EAGAIN {
			break
		}
		if timeout < time.Millisecond {
			break
		}
		timeout -= time.Millisecond
		time.Sleep(time.Millisecond)
	}
	return vfsOS.lockErrorCode(err, IOERR_RDLOCK)
}

func (vfsOSMethods) writeLock(file *os.File, start, len int64, timeout time.Duration) xErrorCode {
	var err error
	for {
		err = unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB)
		if errno, _ := err.(unix.Errno); errno != unix.EAGAIN {
			break
		}
		if timeout < time.Millisecond {
			break
		}
		timeout -= time.Millisecond
		time.Sleep(time.Millisecond)
	}
	return vfsOS.lockErrorCode(err, IOERR_RDLOCK)
}

func (vfsOSMethods) checkLock(file *os.File, start, len int64) (bool, xErrorCode) {
	lock := unix.Flock_t{
		Type:  unix.F_RDLCK,
		Start: start,
		Len:   len,
	}
	if unix.FcntlFlock(file.Fd(), unix.F_GETLK, &lock) != nil {
		return false, IOERR_CHECKRESERVEDLOCK
	}
	return lock.Type != unix.F_UNLCK, _OK
}
