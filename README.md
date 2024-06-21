# ESC Code Generate

using `go generate` to saving time for writing same code in different places.

## domaingen

`domaingen` generates specified file to implement the interface.

### usage

```go
//go:generate domaingen -destination=../../target_file.go -package=targetpkgname (optional) -name=implementedStructName  -replace
type InterfaceYouWantToAutoImplement interface {
    SomeMethod()
}
```

## modelgen

`modelgen` duplicates the structure to another place, and generate the methods to transfer between two structures.

### usage

```go

```
