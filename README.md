# WinGows kernel backup tool

Have u ever experienced **this** and didn't want to launch your favourite tool with admin rights?

![image](https://github.com/user-attachments/assets/62e900cb-0044-4140-aeb9-e27a1dd3af17)

Say no more!

With this tool you can specify needed files via toml config file, or all if needed.

# Requirements
- go

# Effects
![image](https://github.com/user-attachments/assets/822724e4-5672-4ed6-86c7-35593cccb38c)

# Config example
```toml
System32_files_to_dump = [
    "ntoskrnl.exe",
    "win32k.sys",
    "ntdll.dll"
]

# If u want to backup all valid modules from directory, place a single "*"
System32Drivers_files_to_dump = [
    "*"
]

# Otherwise it logs every copied file
Silent = true
```

# Usage
`git clone https://github.com/shv187/winGows_kernel_backuper.git`

`cd winGows_kernel_backuper`

`change config.toml to your needs`

`go run .`

# TODO™
- Rename It
- Restructurize dump layout
- Dump only these files that had changed since the last dump
- [unlikely] Add an option to install it, add it to PATH and dump to some user defined directory instead of relative /dumps for actual QoL

# Notes
This is almost my first Go code, I found the file handling pretty *weird*, lmk if there are any serious issues.
