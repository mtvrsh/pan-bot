version := `git log -1 --pretty=%h`+`test -z "$(git status --porcelain)"|| echo -dirty`

build *ARGS:
    CGO_ENABLED=0 go build -ldflags "-X main.version={{version}} -s -w" {{ARGS}}
alias b := build

test ARGS="-v":
    go test {{ARGS}}
alias t := test

clean:
    rm ./pan-bot
alias c := clean
