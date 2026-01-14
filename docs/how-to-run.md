# How to run RecSIS

- [How to run RecSIS](#how-to-run-recsis)
  - [Clone repository](#clone-repository)
  - [Run](#run)
  - [Summary](#summary)

## Clone repository

**Prerequisites:**
 - Member of RecSIS repo.
 - Set up SSH for GitHub account (see [github docs](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account)).

**Steps:**
```
git clone git@github.com:michalhercik/RecSIS.git
```

## Run

**Prerequisites:**
 - Cloned RecSIS repo (see [Clone](#clone-repository)).
 - Installed Docker (see [Docker docs](https://docs.docker.com/get-docker/)).
 - Installed Go (see [Go docs](https://go.dev/doc/install)).
 - SSH key setup for Acheron (Optional)
    - Being able to access Acheron via SSH using your SSH key with private key
    located at `~/.ssh/id_rsa`.
    - This step allows you to run ELT process which populates RecSIS with data
    from SIS.
    - The requirement can be ignored if you don't mind RecSIS without any SIS
    data.

Before running the RecSIS you need to set environment variables required by
`docker-compose.yml` and webapp. The easiest way is to create a file named
`docker.env` with the required variables and load it in your terminal whenever
you are working with `docker compose`. All `.env` files are not tracked so don't
be afraid of password exposure. Variables needed to be set can be found in
`docker-compose.yml` file under *environment* field of each service.
Alternatively if you run the command `docker compose` it will warn you about
missing variables.

All environment variables with `_PASS` suffix (except `SIS_DB_PASS`) and `MEILI_MASTER_KEY` can be set to any string you want. The string will be used as a password for the corresponding service. Same goes for `POSTGRES_USER` and `POSTGRES_PASSWORD`. `SIS_DB_USER`, `SIS_DB_PASS` and `ACHERON_USER` must be set correctly and if you need them, please contact us at [recsis@email.cz](mailto:recsis@email.cz).

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

**Steps:**
```
docker compose up -d postgres meilisearch elt bert mockcas adminer
```

Now that Meilisearch is running you need to configure it using script. The
script will set aliases, filterable, sortable and searchable attributes.

For **Windows**
```
.\scripts\init_meili.ps1
```

For **Linux**
```
./scripts/init_meili.sh
```

Before running the webapp you need to install [templ](https://templ.guide/) tool
which is responsible for generating HTML templates from `.templ` files
and [wgo](https://github.com/bokwoon95/wgo) which watches live changes in the source files and rebuilds the webapp.

**Steps:**
```
go install github.com/a-h/templ/cmd/templ@v0.2.793
go install github.com/bokwoon95/wgo@latest
```

Lastly you can run the webapp. The best way to do it is using watch script. The
script will automatically rebuild the webapp whenever you change any of the
source files. It also always generates HTML templates.

For **Windows**
```
.\scripts\watch.ps1
```

For **Linux**
```
./scripts/watch.sh
```

If everything went well you should be able to access the webapp at
[https://localhost:8000](https://localhost:8000).

## Summary

For **Windows**:

```
# Clone RecSIS repo
git clone git@github.com:michalhercik/RecSIS.git

# Load environment variables
scripts\init-env.ps1 [.env file path]

# Build & run containers
docker compose up -d postgres meilisearch elt bert mockcas adminer

# Init Meilisearch
.\scripts\init-meili.ps1

# Install templ and wgo
go install github.com/a-h/templ/cmd/templ@v0.2.793
go install github.com/bokwoon95/wgo@latest

# Run webapp
.\scripts\watch.ps1
```

For **Linux**:

```
# Clone RecSIS repo
git clone git@github.com:michalhercik/RecSIS.git

# Load environment variables
source [.env file path]
export $(cut -d= -f1 [.env file path])

# Build & run containers
docker compose up -d postgres meilisearch elt bert mockcas adminer

# Init Meilisearch
./scripts/init-meili.sh

# Install templ and wgo
go install github.com/a-h/templ/cmd/templ@v0.2.793
go install github.com/bokwoon95/wgo@latest

# Run webapp
./scripts/watch.sh
```
