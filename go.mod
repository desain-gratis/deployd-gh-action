module github.com/desain-gratis/deployd-gh-action

go 1.24.0

toolchain go1.24.11

replace github.com/desain-gratis/common => ../common

require (
	github.com/desain-gratis/common v0.0.2-0.20260201145442-11d7f4fcfab6
	github.com/desain-gratis/deployd v0.0.0-20260201150116-affd0bb93c06
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)
