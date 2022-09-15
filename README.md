# miroir

Takes an `io.Reader` and produces two `io.Reader` that will independently consume the upstream reader, buffering data as needed without waiting for the whole input to be consumed first.

This is useful for example when you want to mirror some traffic while affecting the original request as little as possible.

```go
body, mirroredBody := miroir.New(req.Body)

go io.Copy(somewhere, mirroredBody)
process(body)
```

This is conceptually similar to:

```go
body, err := io.ReadAll(req.Body)

go io.Copy(somewhere, bytes.NewReader(body))

process(bytes.NewReader(body))
```

But:

1. it lets you consume either of the returned readers independently of the other, without having to wait until the entire input is finished.
2. it correctly passes upstream errors through both returned readers.

