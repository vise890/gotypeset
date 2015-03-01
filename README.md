# Gotypeset

Tyepesetting as a service. (Multi)markdown comes in, a typeset `.pdf` comes
out...


## Usage

```bash
# in your host
$ vagrant up
$ vagrant ssh
```

```bash
# in the guest vagrant VM
$ export GOPATH=$HOME/go
$ cd ~/go/github.com/vise890/gotypeset
$ go get && go build && ./gotypeset
```

... And visit `http://localhost:8080/` in your browser.

