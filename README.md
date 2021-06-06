[![build](https://github.com/zikwall/go-fileserver/workflows/tests/badge.svg)](https://github.com/zikwall/go-fileserver/actions)

<div align="center">
  <h1>Go-fileserver</h1>
  <h5>Simple, Powerful and Productive file server written in Go</h5>
</div>

### Options

| Name                                | Env                     | 
| :---------------------------------: | :---------------------: |
| `bind-address` (`0.0.0.0:1337`)     |  `BIND_ADDRESS`         |
| `token`                             |  `TOKEN`                |
| `root-file-directory` (`./tmp`)     |  `ROOT_FILE_DIRECTORY`  |
| `enable-secure` (`false`)           |  `ENABLE_SECURE`        |
| `secure-type` (`0`)                 |  `SECURE_TYPE`          |
| `users` (`./tmp`)                   |  `USERS`                |

### Authorization types

- [x] Simple token auth (default)

```shell
--secure-type 0
--enable-secure
--token='token_here'
```

- __token__: With this parameter, you can protect yourself from unauthorized access, it can be empty, if empty-it is generated automatically

From              | Usage 
---               | --- | 
Request header    | `Authorization: Bearer <token>`
Query param       | `?token=<token>` 
Form value        | `token=<token>`

- [x] HTTP Basic auth

```shell
--secure-type 1
--enable-secure
--users='qwx:1337'
--users='qwx2:13372'
--users='qwx3:13373'
```

- [ ] JWT based auth

```shell
// todo
```

### How to use?

#### with single file (basic requests)

```shell
$ curl -i -X POST \
    -H "Content-Type: multipart/form-data" \
    -F "file=@.gitignore" \
    http://localhost:1337?token=123456
```

```shell
$ curl -i -X POST \
    -H "Content-Type: multipart/form-data" \
    -F "file=@.gitignore" \
    http://qwx:1337@localhost:1337
```

```shell
// todo
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

**Simple**
> http://localhost:1337/.gitignore?token=123456

**HTTP Basic**
> http://qwx:1337@localhost:1337/.gitignore

### How to run?

- `$ git clone https://github.com/zikwall/go-fileserver`
- `$ go run . --bind-address localhost:1337 --token 123456 --enable-secure`

#### with Docker

```shell
$ docker run -d -p 1338:1338 \
  -e BIND_ADDRESS='0.0.0.0:1338' \
  -e SECURE_TYPE=1 \
  -e USERS='username:password' \
  -e ENABLE_SECURE='true' \
  -e ENABLE_TSL='true' \
  -e TSL_CERT_FILE='/mnt/ssl.cert' \
  -e TSL_KEY_FILE='/mnt/ssl.key' \
  -v $HOME/tmp:/app/tmp \
  -v $PWD:/mnt \
  --name go-fileserver-example qwx1337/go-fileserver:latest
```

### Tests

- `$ go test -v ./...`
