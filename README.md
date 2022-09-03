# Pan Bot

Very limited discord bot. Do not use.

Discord user interface is in Polish.

## Build

Dependencies: go. Just, git and podman/docker are optional.

* Using [just](https://github.com/casey/just):

```sh
just
```

* Manual build:

```sh
CGO_ENABLED=0 go build
```

## Docker/podman setup:

1. Build docker image:

```sh
just build-dockerfile ${OPTIONAL_REGISTRY}
```

2. Fill `config.json` and create secret:

```sh
podman secret create pan-bot-config config.json
```

3. Create container:

```sh
podman create --secret pan-bot-config --name pan-bot localhost/pan-bot:${VERSION}
```

4. Enable [runit user service](https://docs.voidlinux.org/config/services/user-services.html#per-user-services).

## Misc

* Benchmarks:

```sh
just test -bench-. -v
```
