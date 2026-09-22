module github.com/desain-gratis/deployd-gh-action

go 1.26.0

// replace github.com/desain-gratis/common => ../common

require (
	github.com/desain-gratis/common v0.0.2-0.20260913202850-125f188ddccb
	github.com/desain-gratis/deployd v0.0.16-0.20260922154217-2a22d04af753
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
)
