# Gotypeset

Tyepesetting as a service. (Multi)markdown comes in, a typeset `.pdf` comes
out...


## Usage (for now)

```bash
$ vagrant up
$ vagrant ssh
```

```bash
$ cd /vagrant
$ go build
$ ./gotypeset
```

```
POST http://localhost:8080/typeset {input_file.md}
```
