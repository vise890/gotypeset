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
$ git clone https://github.com/vise890/gotypeset
$ cd gotypeset
$ go build && ./gotypeset
```

... And visit `http://localhost:8080/` in your browser.

