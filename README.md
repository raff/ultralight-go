[![Go Documentation](http://godoc.org/github.com/raff/ultralight-go?status.svg)](http://godoc.org/github.com/raff/ultralight-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/raff/ultralight-go)](https://goreportcard.com/report/github.com/raff/ultralight-go)

# ultralight-go
Go bindings for https://ultralig.ht/

Initial work to implement Go bindiongs for Ultralight (https://github.com/ultralight-ux/ultralight).

For now this requires a few manual steps:

- Get a recent version of the Ultralight SDK. Best option for now is to clone the Ultralight repository and build it
    locally.

- Copy/link the Ultralight SDK in this folder. If you built locally, the SDK is in {repo}/build/SDK.

- Enable setting additional CGO LDFLAGS (at least for MacOS):

      export CGO_LDFLAGS_ALLOW=-Wl,-rpath.*

- Run example:

      go run examples/resize.go

- Run browser:

      cd examples/browser; make 
      ./browser

    The browser needs the HTML assets in examples/browser/assets and expects them in the current directory,
    so you need to "cd" in there.

    It also expects the SDK to be in the current directy, so the Makefile creates a link.
