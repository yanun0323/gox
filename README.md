# ESC Code Generate
使用 `go generate` 來改善大量重工的問題

## esc-model-gen
> 因 `payload/entity/usecase/repository` 在 model 上有大量重複的結構
> 
> `esc-model-gen` 以自動生成 Code 的方式來解決重複撰寫的問題

### Usage
```go
//go:generate esc-model-gen [TAG...]
```

### Layer Dependency
```mermaid
classDiagram
    Payload <|-- Usecase
    Usecase <|-- Repository

    Payload <|-- Entitvy
    Usecase <|-- Entity
    Repository <|-- Entity


    Payload : NewFromUseCase
    Payload : NewFromEntity
    Payload : ToUseCase()
    Payload : ToEntity()

    Usecase : NewFromRepository
    Usecase : NewFromEntity
    Usecase : ToRepository()
    Usecase : ToEntity()

    Repository: NewFromEntity
    Repository: ToEntity()
    Entity: 
```


## esc-domain-gen
gen `usecase/repository` interface 實作使用
