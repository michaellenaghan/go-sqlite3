//go:build (darwin || linux || freebsd || openbsd || netbsd || dragonfly || illumos || sqlite3_flock) && (386 || arm || amd64 || arm64 || riscv64 || ppc64le) && !(sqlite3_noshm || sqlite3_nosys)

package vfs

// SupportsSharedMemory is false on platforms that do not support shared memory.
// To use [WAL without shared-memory], you need to set [EXCLUSIVE locking mode].
//
// [WAL without shared-memory]: https://sqlite.org/wal.html#noshm
// [EXCLUSIVE locking mode]: https://sqlite.org/pragma.html#pragma_locking_mode
const SupportsSharedMemory = true

func (f *vfsFile) SharedMemory() SharedMemory { return f.shm }

// NewSharedMemory returns a shared-memory WAL-index
// backed by a file with the given path.
// It will return nil if shared-memory is not supported,
// or not appropriate for the given flags.
// Only [OPEN_MAIN_DB] databases may need a WAL-index.
// You must ensure all concurrent accesses to a database
// use shared-memory instances created with the same path.
func NewSharedMemory(path string, flags OpenFlag) SharedMemory {
	if flags&OPEN_MAIN_DB == 0 || flags&(OPEN_DELETEONCLOSE|OPEN_MEMORY) != 0 {
		return nil
	}
	return &vfsShm{
		path:     path,
		readOnly: flags&OPEN_READONLY != 0,
	}
}
