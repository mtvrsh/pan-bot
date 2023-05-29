version := `date +%Y%m%d`+"-"+`git log -1 --pretty=%h`+`test -z "$(git status --porcelain)"|| echo -dirty`

build *ARGS:
    CGO_ENABLED=0 go build -ldflags "-X main.version={{version}} -s -w" {{ARGS}}
alias b := build

build-dockerfile *REGISTRY: build
    podman build . -t {{REGISTRY}}pan-bot:{{version}}
alias bd := build-dockerfile

test *ARGS:
    go test {{ARGS}} ./...
alias t := test

clean:
    rm -f ./pan-bot
    podman image rm -f localhost/pan-bot:{{version}}
alias c := clean

push HOST:
    podman save localhost/pan-bot:{{version}} -o /tmp/pan-bot-{{version}}.tar
    scp /tmp/pan-bot-{{version}}.tar {{HOST}}:/tmp/pan-bot-{{version}}.tar
    ssh {{HOST}} podman load -i /tmp/pan-bot-{{version}}.tar
alias p := push
