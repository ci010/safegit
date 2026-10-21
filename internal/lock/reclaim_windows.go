package lock

import (
	"golang.org/x/sys/windows"
	"os"
)

func openReclaim(path string) (*os.File, error) {
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	h, err := windows.CreateFile(p, windows.GENERIC_READ, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE, nil, windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(h), path), nil
}

// Closing the handle releases the whole-file, nonblocking reclamation lock.
func lockReclaim(f *os.File) error {
	return windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, ^uint32(0), ^uint32(0), new(windows.Overlapped))
}
func statIdentity(path string) (os.FileInfo, error) {
	f, e := openReclaim(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return f.Stat()
}
