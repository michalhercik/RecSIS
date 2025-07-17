# RecSIS Development Guide

Before diving into the concrete steps of building and running RecSIS, it would
be beneficial to understand overall structure of the project. If you feel like
ignoring it for now, feel free to skip it and jump right into
[Setup](#setup).

## Introduction

The RecSIS consists of several docker containers and two standalone
applications. The responsibilities can be seen in the diagram below. It's worth
noting that RecSIS is planned to be deployed in production fully dockerized but
because rebuilding webapp as a docker container is much slower the webapp as the
main application is not dockerized as it is rebuildit quite frequently in
development. The second undockerized application (Mock CAS) is a not used in
production at all and the real instance of [CAS](https://cas.cuni.cz/cas/login)
is used instead.

![](../out/docs/services/dev_services.svg)


### Search Engine

Even though we aim to use as little technologies as possible, some are
neccessary to deliver the envisioned UX. Especially since search is a core
feature of RecSIS. After some experiments with PostgreSQL full text search
capabilities We decided to use [Meilisearch](https://www.meilisearch.com/)
because it is easy to set as opposed to for example [Apache
Solr](https://solr.apache.org/), but still provides neccessary features (e.g.:
typo tolerance) as opposed to PostgreSQL. Meilisearch wasn't the only
posibility. [Typesense](https://typesense.org/) was another candidate but we
decided to go with younger Meilisearch because it looked more shiny.

### ELT

Even though implementation of ELT is pretty straightforward and if you are
interested you should read the source code, it is worth mentioning the magic
happining inside the container. ELT is reponsible for extracting, transforming
and loading neccessary data from SIS database into RecSIS database. The catch is
that SIS database is not accessible from outside networks. The container simply
forwards all traffic to a server on which is RecSIS hosted. Which means that to
run ELT you need to have access to server on which RecSIS is hosted. This should
answer possible questions in the future about why ELT needs this or that to
properly execute.

### Adminer

Last Docker container not mentioned is [Adminer](https://www.adminer.org/).
Adminer is a lightweight database management tool which helps to investigate
postgres through simple web interface when developing.

## How to run RecSIS

### Clone repository

**Prerequisities:**
 - Member of RecSIS repo.
 - Setuped SSH for GitHub account (see [github docs](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account)).

**Steps:**
```
git clone git@github.com:michalhercik/RecSIS.git
```

### Run

**Prerequisites:**
 - Cloned RecSIS repo (see [Clone](#clone)).
 - Installed Docker (see [Docker docs](https://docs.docker.com/get-docker/)).
 - Installed Go (see [Go docs](https://go.dev/doc/install)).
 - SSH key setup for Acheron (Optional)
    - Being able to access Acheron via SSH using your SSH key with private key
    located at `~/.ssh/id_rsa`.
    - This step allows you tu run ELT process which populates RecSIS with data
    from SIS.
    - The requirement can be ignored if you don't mind RecSIS witout any SIS
    data.

Before running the RecSIS you need to set environment variables required by
`docker-compose.yml`. The easiest way is to create a file named `docker.env`
with the required variables and load it in your terminal whenever you are
working with `docker compose`. All `.env` files are not tracked so don't be
afraid of password exposure. Variables needed to be set can be found in
`docker-compose.yml` file under *environment* field of each service.
Alternatively if you run the command `docker compose` it will warn you about
missing variables.

You can then load it in your terminal with the following command:

For **Windows**:

```
scripts\init-env.ps1 [.env file path]
```

For **Linux**:

```
source [.env file path]
export $(cut -d= -f1 [.env file path])
```

The next step is to run the `docker compose` command. This will build and run
the necessary containers. The command will also automatically download the
required images if they are not already present on your system.

> NOTE: If you skipped the Acheron SSH setup step you should **not** run the
**elt** service.

**Windows Steps:**
```

```

**Linux Steps:**
```

```

### Summary

<!-- ### Windows

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
``` -->

<!-- ```
# file: docker.env

# SIS database user
DB_USER=
# SIS database password
DB_PASS=
# RecSIS database user
RECSIS_USER=
# RecSIS database password
RECSIS_PASS=
# Acheron user to access SIS database
ACHERON_USER=
# Meilisearch master key
MEILI_MASTER_KEY=
# Postgres admin user
POSTGRES_USER=
# Postgres admin password
POSTGRES_PASSWORD=
# Password for ELT to access RecSIS database
ELT_PASS=
# Password for webapp to access RecSIS database
WEBAPP_PASS=
``` -->