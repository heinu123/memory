# Golang Memory Tool

reading and writing data in the memory of a specified process on a Linux/Android system.


## API

### readAddressDword
```
readAddressDword(pid int, address uintptr) int32
```
> reads a 32-bit integer from the specified process memory address.

Returns:

int32: The value corresponding to the address

### readAddressFloat
```
readAddressFloat(pid int, address uintptr) float32
```
> reads a 32-bit float from the specified process memory address.

Returns:

float32: The value corresponding to the address

### writeAddressDword
```
writeAddressDword(pid int, address uintptr, value int32) error
```
> writes a 32-bit integer to the specified process memory address.


### writeAddressFloat
```
writeAddressFloat(pid int, address uintptr, value float32) error
```
> writes a 32-bit float to the specified process memory address.

### getModuleBase
```
getModuleBase(pid int, moduleName string) uint64
```
> retrieves the base address of the specified module in the process memory.

Returns:

uint64: The base address of the specified module in the process memory.


### getPID
```
getPID(packageName string) int
```
> Get the pid of the specified process (such as package name/process name)

Returns:

int: PID




## Usage

```go
import "github.com/heinu123/memory"
```
```go
    pid := memory.GetPID("exampleProcess")
    baseAddress := memory.GetModuleBase(pid, "libUE4.so")
    
    // Read a DWORD from the memory
    value := memory.ReadAddressDword(pid, uintptr(baseAddress))
    
    // Write a DWORD to the memory
    err := memory.WriteAddressDword(pid, uintptr(baseAddress), 42)
    if err != nil {
        // Handle the error
    }
    

```
Feel free to customize the package usage according to your specific requirements.




Is this conversation helpful so far?


Send a message

