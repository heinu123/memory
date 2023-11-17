package memory

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "math"
    "unsafe"
    "syscall"
    "encoding/binary"
)

func ReadAddressDword(pid int, address uintptr) (int32) {
    memPath := fmt.Sprintf("/proc/%d/mem", pid)
    handle, err := syscall.Open(memPath, syscall.O_RDONLY, 0)
    if err != nil {
        return 0
    }
    defer syscall.Close(handle)

    data := make([]byte, 4)
    _, _, errno := syscall.Syscall6(syscall.SYS_PREAD64, uintptr(handle), uintptr(unsafe.Pointer(&data[0])), 4, address, 0, 0)
    if errno != 0 {
        return 0
    }

    return int32(binary.LittleEndian.Uint32(data))
}

func ReadAddressFloat(pid int, address uintptr) (float32) {
    memPath := fmt.Sprintf("/proc/%d/mem", pid)
    handle, err := syscall.Open(memPath, syscall.O_RDONLY, 0)
    if err != nil {
        return 0
    }
    defer syscall.Close(handle)

    data := make([]byte, 4)
    _, _, errno := syscall.Syscall6(syscall.SYS_PREAD64, uintptr(handle), uintptr(unsafe.Pointer(&data[0])), 4, address, 0, 0)
    if errno != 0 {
        return 0
    }

    bits := binary.LittleEndian.Uint32(data)
    return math.Float32frombits(bits)
}

func WriteAddressDword(pid int, address uintptr, value int32) error {
    memPath := fmt.Sprintf("/proc/%d/mem", pid)
    handle, err := syscall.Open(memPath, syscall.O_RDWR, 0)
    if err != nil {
        return err
    }
    defer syscall.Close(handle)

    data := make([]byte, 4)
    binary.LittleEndian.PutUint32(data, uint32(value))

    _, _, errno := syscall.Syscall6(syscall.SYS_PWRITE64, uintptr(handle), uintptr(unsafe.Pointer(&data[0])), 4, address, 0, 0)
    if errno != 0 {
        return errno
    }

    return nil
}

func WriteAddressFloat(pid int, address uintptr, value float32) error {
    memPath := fmt.Sprintf("/proc/%d/mem", pid)
    handle, err := syscall.Open(memPath, syscall.O_RDWR, 0)
    if err != nil {
        return err
    }
    defer syscall.Close(handle)

    bits := math.Float32bits(value)
    data := make([]byte, 4)
    binary.LittleEndian.PutUint32(data, bits)

    _, _, errno := syscall.Syscall6(syscall.SYS_PWRITE64, uintptr(handle), uintptr(unsafe.Pointer(&data[0])), 4, address, 0, 0)
    if errno != 0 {
        return errno
    }

    return nil
}

func GetModuleBase(pid int, moduleName string) uint64 {
    filename := fmt.Sprintf("/proc/%d/maps", pid)
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return 0
    }

    lines := strings.Split(string(content), "\n")
    for _, line := range lines {
        if strings.Contains(line, moduleName) {
            fields := strings.Fields(line)
            addrStr := strings.Split(fields[0], "-")[0]
            addr, err := strconv.ParseUint(addrStr, 16, 64)
            if err != nil {
                return 0
            }
            if addr == 0x8000 {
                addr = 0
            }
            return addr
        }
    }

    return 0
}

func GetPID(packageName string) int {
    dir, err := os.Open("/proc")
    if err != nil {
        return 0
    }
    defer dir.Close()

    entries, err := dir.Readdirnames(-1)
    if err != nil {
        return 0
    }

    for _, entry := range entries {
        if entry == "." || entry == ".." {
            continue
        }

        info, err := os.Lstat(filepath.Join("/proc", entry))
        if err != nil || !info.IsDir() {
            continue
        }

        cmdlinePath := filepath.Join("/proc", entry, "cmdline")
        cmdlineBytes, err := ioutil.ReadFile(cmdlinePath)
        if err != nil {
            continue
        }

        cmdline := string(cmdlineBytes)
        if strings.Contains(cmdline, packageName) {
            pid, err := strconv.Atoi(entry)
            if err == nil {
                return pid
            }
        }
    }

    return 0
}