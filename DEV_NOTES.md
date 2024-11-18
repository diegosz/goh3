# Development Notes

> [!WARNING]
> This is a work in progress, it's not working yet.

## Update `ch3` Package

Update the `ch3` package with the latest changes from the [`h3`](https://github.com/diegosz/h3) package.

```shell
rm .drone.yml
rm ch3/a_darwin_arm64.go
rm ch3/a_linux_arm64.go
rm ch3/capi_darwin_arm64.go
rm ch3/capi_linux_arm64.go
cp <h3-directory>/a_linux_amd64.go ch3/a_linux_amd64.go
cp <h3-directory>/capi_linux_amd64.go ch3/capi_linux_amd64.go
cp <h3-directory>/helper.go ch3/helper.go
cp /home/diegos/_dev/github/diegosz/h3/a_linux_amd64.go ch3/a_linux_amd64.go
cp /home/diegos/_dev/github/diegosz/h3/capi_linux_amd64.go ch3/capi_linux_amd64.go
cp /home/diegos/_dev/github/diegosz/h3/helper.go ch3/helper.go
go get -u ./...
go mod tidy
```

- [Function name changes | H3](https://h3geo.org/docs/library/migration-3.x/functions/#general-function-names)

## FIXME

Running test in the IDE are not working as expected. It works fine in debug mode, but fails to run in testing.

Calling library functions returns empty. Something smells...

May be something related to the [thread-local storage](https://groups.google.com/g/golang-nuts/c/tGamryo50BY).

¯\\_(ツ)_/¯

### Resources

- <https://groups.google.com/g/golang-nuts/c/tGamryo50BY>
- <https://stackoverflow.com/questions/43292365/calling-functions-inside-a-lockosthread-goroutine>
- <https://stackoverflow.com/questions/25361831/benefits-of-runtime-lockosthread-in-golang>
- <https://github.com/golang/go/issues/21827>
