# Gox Code Generator

using `go generate` to saving time for writing same code in different places.

## domaingen

`domaingen` generates specified file to implement the interface.

### usage

```bash
-h                              show usage
-name                           implemented structname                 -name=usecase
-package        (require)       implemented struct package name
-destination    (require)       generated filepath                     -destination=../../usecase/member_usecase.go
-replace                        force replace exist struct/funcmethod
-constructor                    generate constructor function
example:
//go:generate domaingen -destination=../../usecase/member.go-name=usecase -replace -constructor
```

```go
//go:generate domaingen -destination=../../target_file.go -package=targetpkgname (optional) -name=implementedStructName  -replace -constructor
type InterfaceYouWantToAutoImplement interface {
    SomeMethod()
}
```

## modelgen

#### coming soon...
