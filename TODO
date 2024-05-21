##  Adding support for Powershell


### From [Issue #140: "Windows PowerShell Support" in the original repo (knqyf263/pet)](https://github.com/knqyf263/pet/issues/140)

> In Windows PowerShell, creating command entries works as expected; however the `pet exec` function defaults to running commands in cmd.exe and not PowerShell.  Can `pet` be configured to default to pwsh.exe instead?

Currently the source contains a file named [`util_windows.go`](https://github.com/poa00/pet/blob/main/cmd/util_windows.go) (as well as `util_unix.go`).  These define the `run` function which is used to execute the command in question.  Alternative `run` that uses PowerShell needs to be defined.

#### Potential Solutions:

1.  Change the `run` signature to also accept the 'shell' type, possibly by detecting if the `pet` is currently running in PowerShell).
 
2.  Create a global config flag that decides which shell to use, effectively limiting the user to using either PowerShell or Command Console exclusively, but not both simultaneously.
 
**Solution 2** is likely easier to implement and less disruptive. It is uncertain, but the 'selection' _may_ require updating as as well.
