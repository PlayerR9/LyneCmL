# table
A Go package used for generating tables implemented without bounds. It also features some already generated tables.


## Table of Contents

1. [Table of Contents](#table-of-contents)
2. [Tool](#tool)
   - [Installation](#installation)
   - [Usage](#usage)
3. [Documentation](#documentation)
4. [Content](#content)


## Tool

### Installation

To install the tool, run the following command:
```
go get -u github.com/PlayerR9/table/cmd
```


### Usage

Once imported, you can use the tool to generate tables for your own types. Like so:
```go
import _ "github.com/PlayerR9/table"

// go:generate go run table/cmd -name=Foo -type=int
```

This generates a table with the name "Foo" whose cells are of type "int".

The type generated will be in the same package as the tool. Make sure to read the documentation of the tool before using it.


## Documentation

```markdown
This command generates a table of the given type implemented as a boundless table (i.e., out-of-bounds errors are not thrown).

To use it, run the following command:

   go:generate go run table/cmd -name=<type_name> -type=<type> [ -g=<generics>] [ -o=<output_file> ]

**Flag: Name**

The "name" flag is used to specify the name of the table. As such, it must be set and,
not only does it have to be a valid Go identifier, but it also must start with an upper case letter.

**Flag: Type**

The "fields" flag is used to specify type of the table's cells. This flag must be set.

For instance, running the following command:

   go:generate table -name=Table -type=int

will generate a table with the following structure:

   type Table struct {
      table [][]int
   }

It is important to note that spaces are not allowed in any of the flags.

Also, it is possible to specify generics by following the value with the generics between square brackets;
like so: "MyType[T,C]"


**Flag: Generics**

This optional flag is used to specify the type(s) of the generics. However, this only applies if at least one
generic type is specified in the type flag. If none, then this flag is ignored.

As an edge case, if this flag is not specified but the type flag contains generics, then
all generics are set to the default value of "any".

Its argument is specified as a list of key-value pairs where each pair is separated
by a comma (",") and a slash ("/") is used to separate the key and the value. The key indicates the name of
the generic and the value indicates the type of the generic.

For instance, running the following command:

   go:generate table -name=Table -type=MyType[T] -g=T/any

will generate a table with the following fields:

   type Table[T any] struct {
      table [][]T
   }


**Flag: Output File**

This optional flag is used to specify the output file. If not specified, the output will be written to
standard output, that is, the file "<type_name>_table.go" in the root of the current directory.
```


## Content

Here are all the pregenerated files:
- [bool](bool.go)
- [byte](byte.go)
- [int](int.go)
- [int8](int8.go)
- [int16](int16.go)
- [int32](int32.go)
- [int64](int64.go)
- [float32](float32.go)
- [float64](float64.go)
- [rune](rune.go)
- [string](string.go)
- [uint](uint.go)
- [uint8](uint8.go)
- [uint16](uint16.go)
- [uint32](uint32.go)
- [uint64](uint64.go)
- [uintptr](uintptr.go)
- [error](error.go)
- [complex128](complex128.go)
- [complex64](complex64.go)
- [generic](generic.go)