# MessagePack encoding for Golang

[![Build Status](https://travis-ci.org/vmihailenco/msgpack.svg)](https://travis-ci.org/vmihailenco/msgpack)
[![PkgGoDev](https://pkg.go.dev/badge/code.byted.org/ad/msgpack_extstr)](https://pkg.go.dev/code.byted.org/ad/msgpack_extstr)
[![Documentation](https://img.shields.io/badge/msgpack-documentation-informational)](https://msgpack.uptrace.dev/)
[![Chat](https://discordapp.com/api/guilds/752070105847955518/widget.png)](https://discord.gg/rWtp5Aj)

> :heart:
> [**Uptrace.dev** - All-in-one tool to optimize performance and monitor errors & logs](https://uptrace.dev/?utm_source=gh-msgpack&utm_campaign=gh-msgpack-var2)

- Join [Discord](https://discord.gg/rWtp5Aj) to ask questions.
- [Documentation](https://msgpack.uptrace.dev)
- [Reference](https://pkg.go.dev/code.byted.org/ad/msgpack_extstr)
- [Examples](https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#pkg-examples)

Other projects you may like:

- [Bun](https://bun.uptrace.dev) - fast and simple SQL client for PostgreSQL, MySQL, and SQLite.
- [treemux](https://github.com/myhyh/treemux) - high-speed, flexible, tree-based HTTP router
  for Go.

## Features

- Primitives, arrays, maps, structs, time.Time and interface{}.
- Appengine \*datastore.Key and datastore.Cursor.
- [CustomEncoder]/[CustomDecoder] interfaces for custom encoding.
- [Extensions](https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#example-RegisterExt) to encode
  type information.
- Renaming fields via `msgpack:"my_field_name"` and alias via `msgpack:"alias:another_name"`.
- Omitting individual empty fields via `msgpack:",omitempty"` tag or all
  [empty fields in a struct](https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#example-Marshal-OmitEmpty).
- [Map keys sorting](https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#Encoder.SetSortMapKeys).
- Encoding/decoding all
  [structs as arrays](https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#Encoder.UseArrayEncodedStructs)
  or
  [individual structs](https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#example-Marshal-AsArray).
- [Encoder.SetCustomStructTag] with [Decoder.SetCustomStructTag] can turn msgpack into drop-in
  replacement for any tag.
- Simple but very fast and efficient
  [queries](https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#example-Decoder.Query).

[customencoder]: https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#CustomEncoder
[customdecoder]: https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#CustomDecoder
[encoder.setcustomstructtag]:
  https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#Encoder.SetCustomStructTag
[decoder.setcustomstructtag]:
  https://pkg.go.dev/code.byted.org/ad/msgpack_extstr#Decoder.SetCustomStructTag

## Installation

msgpack supports 2 last Go versions and requires support for
[Go modules](https://github.com/golang/go/wiki/Modules). So make sure to initialize a Go module:

```shell
go mod init github.com/my/repo
```

And then install msgpack/v5 (note _v5_ in the import; omitting it is a popular mistake):

```shell
go get code.byted.org/ad/msgpack_extstr
```

## Quickstart

```go
import "code.byted.org/ad/msgpack_extstr"

func ExampleMarshal() {
    type Item struct {
        Foo string
    }

    b, err := msgpack.Marshal(&Item{Foo: "bar"})
    if err != nil {
        panic(err)
    }

    var item Item
    err = msgpack.Unmarshal(b, &item)
    if err != nil {
        panic(err)
    }
    fmt.Println(item.Foo)
    // Output: bar
}
```

## See also

- [Fast and flexible ORM for sql.DB](https://bun.uptrace.dev)
