# go-2048-wasm 

Yet another WebAssembly demo.
 
![](2048.png)
# Build & Run 

```
$ GOARCH=wasm GOOS=js go build -o 2048.wasm 2048.go 
$ goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'
```
