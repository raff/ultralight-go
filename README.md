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
