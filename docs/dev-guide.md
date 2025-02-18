# RecSIS Development Guide

**LAST UPDATED**: 18.2. 2025

- [RecSIS Development Guide](#recsis-development-guide)
  - [External Dependencies (summary)](#external-dependencies-summary)
  - [Clone](#clone)
  - [Build / Run](#build--run)
  - [Live Reload](#live-reload)


## External Dependencies (summary)
  - [Go](https://go.dev/)
  - [Templ](https://templ.guide/) - Compiled HTML templates
  - [Docker](https://www.docker.com/)
  - [wgo](github.com/bokwoon95/wgo) - Live reload

## Clone

**Prerequisities:**
 - Member of RecSIS repo.
 - Setuped SSH for account (see
[github docs](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account)).

**Steps:**
```
git clone git@github.com:michalhercik/RecSIS.git
git switch dev
```

## Build / Run

**Prerequisites:**
 - Cloned RecSIS repo (see [Clone](#clone)).
 - Installed Go (see [Go docs](https://go.dev/doc/install)).
    - Linux: add Go bin to path `export PATH="$HOME/go/bin:$PATH"`
    - Windows: **TODO**

### Windows

#### Build

**Steps:**
```
go mod download
go install github.com/a-h/templ/cmd/templ@v0.2.793
./scripts/build.ps1
```
#### Run

**Steps:**
```
go mod download
go install github.com/a-h/templ/cmd/templ@v0.2.793
docker compose up --build -d
./scripts/init_meili.ps1
./scripts/run.ps1
```

### Linux

#### Build

**Steps:**
```
go mod download
go install github.com/a-h/templ/cmd/templ@v0.2.793
./scripts/build.sh
```

#### Run

**Steps:**
```
go mod download
go install github.com/a-h/templ/cmd/templ@v0.2.793
docker compose up --build -d
./scripts/init_meili.sh
./scripts/run.sh
```

## Live Reload

**Prerequisites:**
 - Cloned RecSIS repo (see [Clone](#clone)).
 - Installed Go (see [Go docs](https://go.dev/doc/install)).
    - Linux: add Go bin to path `export PATH="$HOME/go/bin:$PATH"`
    - Windows: **TODO**

### Windows

**Steps:**
```
go mod download
go install github.com/bokwoon95/wgo@v0.5.11
go install github.com/a-h/templ/cmd/templ@v0.2.793
docker compose up --build -d
./scripts/init_meili.ps1
./scripts/watch.ps1
```

### Linux

**Steps:**
```
go mod download
go install github.com/bokwoon95/wgo@v0.5.11
go install github.com/a-h/templ/cmd/templ@v0.2.793
docker compose up --build -d
./scripts/init_meili.sh
./scripts/watch.sh
```