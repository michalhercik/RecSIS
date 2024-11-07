# RecSIS Development Guide

- [RecSIS Development Guide](#recsis-development-guide)
  - [Clone](#clone)
  - [Build](#build)
  - [Live Reload](#live-reload)


## Clone

**Prerequisities:**
 - Member of RecSIS repo.
 - Setuped SSH for account (see
[github docs](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account)).

```
git clone git@github.com:michalhercik/RecSIS.git
git switch dev
```

## Build

**Prerequisites:**
 - Cloned RecSIS repo (see [Clone](#clone)).
 - Installed Go (see [Go docs](https://go.dev/doc/install)).

```
go mod download
templ generate
go build .
```

## Live Reload

**Prerequisites:**
 - Cloned RecSIS repo (see [Clone](#clone)).
 - Installed Go (see [Go docs](https://go.dev/doc/install)).

### Windows

```
go mod download
go install github.com/bokwoon95/wgo@latest
./watch.ps1
```

### Other

```
go mod download
go install github.com/bokwoon95/wgo@latest
wgo -file=".go" -file=".templ" -xfile="_templ.go" templ generate :: go run .
```