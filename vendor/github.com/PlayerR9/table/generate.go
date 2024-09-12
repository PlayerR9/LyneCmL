//go:generate go run cmd/main.go -name=Table -type=T -g=T/any -o=generic.go
//go:generate go run cmd/main.go -name=BoolTable -type=bool -o=bool.go
//go:generate go run cmd/main.go -name=ByteTable -type=byte -o=byte.go
//go:generate go run cmd/main.go -name=Complex64Table -type=complex64 -o=complex64.go
//go:generate go run cmd/main.go -name=Complex128Table -type=complex128 -o=complex128.go
//go:generate go run cmd/main.go -name=ErrorTable -type=error -o=error.go
//go:generate go run cmd/main.go -name=Float32Table -type=float32 -o=float32.go
//go:generate go run cmd/main.go -name=Float64Table -type=float64 -o=float64.go
//go:generate go run cmd/main.go -name=IntTable -type=int -o=int.go
//go:generate go run cmd/main.go -name=Int8Table -type=int8 -o=int8.go
//go:generate go run cmd/main.go -name=Int16Table -type=int16 -o=int16.go
//go:generate go run cmd/main.go -name=Int32Table -type=int32 -o=int32.go
//go:generate go run cmd/main.go -name=Int64Table -type=int64 -o=int64.go
//go:generate go run cmd/main.go -name=RuneTable -type=rune -o=rune.go
//go:generate go run cmd/main.go -name=StringTable -type=string -o=string.go
//go:generate go run cmd/main.go -name=UintTable -type=uint -o=uint.go
//go:generate go run cmd/main.go -name=Uint8Table -type=uint8 -o=uint8.go
//go:generate go run cmd/main.go -name=Uint16Table -type=uint16 -o=uint16.go
//go:generate go run cmd/main.go -name=Uint32Table -type=uint32 -o=uint32.go
//go:generate go run cmd/main.go -name=Uint64Table -type=uint64 -o=uint64.go
//go:generate go run cmd/main.go -name=UintptrTable -type=uintptr -o=uintptr.go

package table
