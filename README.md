# diskio
diskio is a tool to perform high load disk I/O for test purpose.

# Install

```shell script
go install github.com/devlights/diskio
```

# Run

## windows

```shell script
diskio.exe
```

or

```shell script
diskio.exe -g number_of_concurrent_proc(default 100) -b block_size(default 1024)
```

## MacOS and Linux

```shell script
diskio
```

or

```shell script
diskio -g number_of_concurrent_proc(default 100) -b block_size(default 1024)
```
