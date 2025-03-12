# Build

## Dockerfile

- build (linux/arm64)

    ```shell
    docker buildx build --platform linux/arm64 -t couchbase_bridge .
    ```

---

# Running on Local Environment

## Couchbase Port-forward from EKS to Local device

- eks에서 동작 중인 couchbase pod을 port-forward하여 local device에서 web ui에 접속할 수 있습니다.
  `http://localhost:8091/`

  ```shell
  kubectl port-forward svc/couchbase -n dclo-cspm 8091:8091
  ```

- `couchbase_bridge`를 사용하기 위해 필요한 port를 모두 port-forward 하는 경우
  ```shell
  kubectl port-forward svc/couchbase -n dclo-cspm 8091:8091 8092:8092 8093:8093 8094:8094 11210:11210
  ```

- couchbase에서 사용하는 ports 중 일부를 추리면 다음과 같습니다.
    - 8091
        - Cluster administration REST/HTTP traffic, including Couchbase Web Console
    - 8092
        - Views and XDCR access
    - 8093
        - Query service REST/HTTP traffic
    - 8094
        - Search Service REST/HTTP traffic

- [Couchbase ports document](https://docs.couchbase.com/server/current/install/install-ports.html)


