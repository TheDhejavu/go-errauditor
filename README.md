# go-erraudior

An interesting attempt in building an error audior in golang, this auditor scans through your codebase and returns the line number of the function , the name of the function and the list of errors that are returned within the function.

## Installation

### By go get

```
go get github.com/thedhejavu/errauditor/cmd/errauditor
```

## Usage

```bash
errauditor ./...
```
