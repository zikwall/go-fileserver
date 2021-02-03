<div align="center">
  <h1>Go-fileserver</h1>
  <h5>Simple, Powerful and Productive file server written in Go</h5>
</div>

### Options

- __token__: With this parameter, you can protect yourself from unauthorized access, it can be empty, if empty-it is generated automatically
  - Request header `Authorization: Bearer <token>`  
  - Query param `?token=<token>`
  - Form value `token=<token>`
  
  > The list is arranged in order of priority of token processing

- __enable-secure__: Enabling and disabling access token protection - `true` or `false`, default `false`

### How to use?

#### with single file

```shell
$ curl -i -X POST \
    -H "Content-Type: multipart/form-data" \
    -F "file=@.gitignore" \
    http://localhost:1337?token=123456
```

#### with multiple files

```shell
$ curl -i -X POST \
    -H "Content-Type: multipart/form-data" \
    -F "files[]=@.gitignore" \
    -F "files[]=@README.md" \
    http://localhost:1337?token=123456
```

#### with data binary format

```shell
$ curl -i -X POST \
    --data-binary @.gitignore \
    http://localhost:1337/.gitignore?token=123456
```
> The name from the request URI will be taken into account as the file name

### How to view files

To view the files just click on the link with the domain, and the file name at the end: 

> http://localhost:1337/.gitignore?token=123456

### How to run?

- `$ git clone https://github.com/zikwall/go-fileserver`
- `$ go run . --bind-address localhost:1337 --token 123456 --enable-secure`

#### with Docker

```shell
$ docker run -d -p 1337:1337 \
    -v $HOME/tmp:/app/tmp \
    -e BIND_ADDRESS='0.0.0.0:1337' \
    -e TOKEN='123456' \
    --name go-fileserver-example qwx1337/go-fileserver:latest
```

### Tests

- `$ go test -v ./...`