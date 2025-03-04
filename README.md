# Build

### 설치형

```
SET GOOS=linux
SET GOARCH=amd64
go build -o go_couchbase
```

### SaaS

```
SET GOOS=linux
SET GOARCH=arm
go build -o go_couchbase
```

### Docker

- build (linux/arm64)

    ```shell
    docker buildx build --platform linux/arm64 -t couchbase_bridge .
    ```

# Application Name

couchbase_bridge

## Troubleshoot

eks에서 동작 중인 couchbase pod을 port-forward하여 local device에서 web ui에 접속할 수 있습니다.

```shell
kubectl port-forward svc/couchbase -n dclo-cspm 8091:8091
```

`http://localhost:8091/`