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
$ cd /vagrant
$ go run ./gotypeset.go
```

From your host:
```
POST http://localhost:8080/typeset {input_file.md}
```
