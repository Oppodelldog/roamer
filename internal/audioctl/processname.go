package audioctl

import (
	"errors"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Th32csSnapProcess is described in https://msdn.microsoft.com/de-de/library/windows/desktop/ms682489(v=vs.85).aspx
const Th32csSnapProcess = 0x00000002

func ProcessName(pid int) string {
	ps, err := processes()
	if err != nil {
		log.Fatal(err)
	}

	for _, proc := range ps {
		if proc.ProcessID == pid {
			return proc.Exe
		}
	}

	return ""
}

// WindowsProcess is an implementation of Process for Windows.
type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

func processes() ([]WindowsProcess, error) {
	var (
		entry  windows.ProcessEntry32
		handle windows.Handle
		err    error
	)

	handle, err = windows.CreateToolhelp32Snapshot(Th32csSnapProcess, 0)

	if err != nil {
		return nil, err
	}

	defer func() { must(windows.CloseHandle(handle)) }()

	entry.Size = uint32(unsafe.Sizeof(entry))

	err = windows.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)

	for {
		results = append(results, newWindowsProcess(&entry))

		err = windows.Process32Next(handle, &entry)
		if err != nil {
			if errors.Is(err, syscall.ERROR_NO_MORE_FILES) {
				return results, nil
			}

			return nil, err
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func newWindowsProcess(e *windows.ProcessEntry32) WindowsProcess {
	var end int

	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}
