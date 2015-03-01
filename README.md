# Gotypeset

Tyepesetting as a service. (Multi)markdown comes in, a typeset `.pdf` comes
out...


## Usage

```bash
# in your host
$ git clone https://github.com/vise890/gotypeset
$ cd gotypeset
$ vagrant up
$ vagrant ssh
```

```bash
# in the guest vagrant VM
$ cd
$ echo 'export GOPATH=$HOME/go' >> .bashrc
$ source .bashrc
$ cd ~/go/github.com/vise890/gotypeset
$ go get && go build && ./gotypeset
```

... And visit `http://localhost:9000/` in your browser.
