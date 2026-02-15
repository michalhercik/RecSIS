# How to run RecSIS

- [How to run RecSIS](#how-to-run-recsis)
  - [Clone repository](#clone-repository)
  - [Run](#run)
  - [Summary](#summary)

## Clone repository

The project provides two tagged releases:

- **v0-alpha** – version before the Degree Plans page update  
- **v1-alpha** – version after the update

> The `main` branch may change in the future.  
> To run a specific version of the system, always checkout a release tag instead of relying on `main`.

**Prerequisites:**
- Configure SSH access to GitHub (see [GitHub documentation](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account))

**Steps:**

Clone the repository:

```bash
git clone git@github.com:michalhercik/RecSIS.git
cd RecSIS
```

List available releases:

```bash
git tag
```

Checkout a specific version:

```bash
# Version BEFORE degree plans page update
git checkout v0-alpha

# Version WITH new degree plans page
git checkout v1-alpha
```

> On newer Git versions you can also use:
>
> ```bash
> git switch v1-alpha
> ```

Alternatively, if you already cloned the repository in the past, your main branch is probably up to date with v0-alpha. If you want to run v1-alpha, you have to pull from origin and then you can just switch to the tag as shown above. If you want to run v0-alpha, you can just stay on the main branch, but make sure it is up to date with the latest changes:

```bash
git pull origin main
git checkout v1-alpha
```

## Run

**Prerequisites:**
 - Cloned RecSIS repo (see [Clone](#clone-repository)).
 - Installed Docker (see [Docker docs](https://docs.docker.com/get-docker/)).
 - Installed Go (see [Go docs](https://go.dev/doc/install)).
 - SSH key setup for Acheron (Optional)
    - Being able to access Acheron via SSH using your SSH key with private key
    located at `~/.ssh/id_rsa`.
    - This step allows you tu run ELT process which populates RecSIS with data
    from SIS.
    - The requirement can be ignored if you don't mind RecSIS without any SIS
    data.

Before running the RecSIS you need to set environment variables required by `docker-compose.yml` and webapp. The easiest way is to create a file named `docker.env` with the required variables and load it in your terminal whenever you are working with `docker compose` or the webapp in general. All `.env` files are not tracked so don't be afraid of password exposure. Variables needed to be set can be found in `docker-compose.yml` file under *environment* field of each service. Alternatively if you run the command `docker compose` it will warn you about missing variables.

All environment variables with `_PASS` suffix (except `SIS_DB_PASS`) and `MEILI_MASTER_KEY` can be set to any string you want. The string will be used as a password for the corresponding service. Same goes for `POSTGRES_USER` and `POSTGRES_PASSWORD`. `POSTGRES_OWNER` must be set to `postgres`. `SIS_DB_USER`, `SIS_DB_PASS` and `ACHERON_USER` must be set correctly and if you need them, please contact us at [recsis@email.cz](mailto:recsis@email.cz). For ELT you might want to also set `ELT_MAX_OPEN_CONNS`, `ELT_MAX_IDLE_CONNS`, and `ELT_CONN_MAX_LIFETIME`, first two should be numbers, and the last one should be a duration string parsable by Go's `time.ParseDuration()` function. Although, we have fallback values for these three variables in case they are not set, so you can ignore them if you don't know what to set them to or you don't need ELT process.

You can then load it in your terminal with the following command:

For **Windows**:

```shell
scripts\init-env.ps1 [.env file path]
```

For **Linux**:

```bash
source [.env file path]
export $(cut -d= -f1 [.env file path])
```

The next step depends on which version you started and which version you want to run. If you already have v0-alpha running and you want to run v0-alpha, do nothing.

If you started with v0-alpha and you want to run v1-alpha, you need to migrate the database to v1-alpha version. Then you need to rebuild the elt container to have data compatible with v1-alpha. You can do it with the following commands:

> NOTE: If you skipped the Acheron SSH setup step you should **not** run the **elt** service.

**Steps:**

For **Windows**
```shell
# migrate tables populated by ELT
.\scripts\migrate.ps1 -MigrationsDir .\migrates\elt-tables -Container recsis-postgres
# Migrate user tables (without losing user data)
.\scripts\migrate.ps1 -MigrationsDir .\migrates\v1-alpha -Container recsis-postgres

docker compose up --build elt
```

For **Linux**
```bash
# migrate tables populated by ELT
./scripts/migrate.sh ./migrates/elt-tables recsis-postgres
# Migrate user tables (without losing user data)
./scripts/migrate.sh ./migrates/v1-alpha recsis-postgres

docker compose up --build elt
```

If you cloned the repository for the first time and you want to run v0-alpha, you can just run the command `docker compose` as shown in the next step. This will build and run the necessary containers. The command will also automatically download the required images if they are not already present on your system.

> NOTE: If you skipped the Acheron SSH setup step you should **not** run the **elt** service.

**Steps:**
```bash
docker compose up -d postgres meilisearch elt bert mockcas adminer
```

If you want to run v1-alpha, you must switch to v1-alpha version as shown in the [Clone](#clone-repository) step and then run the command `docker compose` as shown above but without the elt service.

Then you have to migrate to v1-alpha version as shown in the previous variant. Finally you have to build and run the elt service to populate the database with data from SIS.

> NOTE: If you skipped the Acheron SSH setup step you should **not** run the **elt** service.

**Steps:**
```bash
docker compose up -d postgres meilisearch bert mockcas adminer
```

For **Windows**
```shell
.\scripts\migrate.ps1 -MigrationsDir .\migrates\elt-tables -Container recsis-postgres
.\scripts\migrate.ps1 -MigrationsDir .\migrates\v1-alpha -Container recsis-postgres
```

For **Linux**
```bash
./scripts/migrate.sh ./migrates/elt-tables recsis-postgres
./scripts/migrate.sh ./migrates/v1-alpha recsis-postgres
```

```bash
docker compose up --build elt
```

For any version, now that Meilisearch is running you need to configure it using script. The script will set aliases, filterable, sortable and searchable attributes.

For **Windows**
```shell
.\scripts\init-meili.ps1
```

For **Linux**
```bash
./scripts/init-meili.sh
```

Before running the webapp you need to install [templ](https://templ.guide/) tool
which is responsible for generating HTML templates from `.templ` files
and [wgo](https://github.com/bokwoon95/wgo) which watches live changes in the source files and rebuilds the webapp.

**Steps:**
```bash
go install github.com/a-h/templ/cmd/templ@v0.2.793
go install github.com/bokwoon95/wgo@latest
```

Lastly you can run the webapp. The best way to do it is using watch script. The
script will automatically rebuild the webapp whenever you change any of the
source files. It also always generates HTML templates.

For **Windows**
```shell
.\scripts\watch.ps1
```

For **Linux**
```bash
./scripts/watch.sh
```

If everything went well you should be able to access the webapp at
[https://localhost:8000](https://localhost:8000).

## Summary

For **Windows**:

```shell
# Clone RecSIS repo
git clone git@github.com:michalhercik/RecSIS.git

# Switch to the desired version
git switch v1-alpha

# Load environment variables
scripts\init-env.ps1 [.env file path]

# Build & run containers
docker compose up -d postgres meilisearch bert mockcas adminer

# Migrate database to v1-alpha version
.\scripts\migrate.ps1 -MigrationsDir .\migrates\elt-tables -Container recsis-postgres
.\scripts\migrate.ps1 -MigrationsDir .\migrates\v1-alpha -Container recsis-postgres

# Populate database with data from SIS
docker compose up --build elt

# Init Meilisearch
.\scripts\init-meili.ps1

# Install templ and wgo
go install github.com/a-h/templ/cmd/templ@v0.2.793
go install github.com/bokwoon95/wgo@latest

# Run webapp
.\scripts\watch.ps1
```

For **Linux**:

```bash
# Clone RecSIS repo
git clone git@github.com:michalhercik/RecSIS.git

# Switch to the desired version
git switch v1-alpha

# Load environment variables
source [.env file path]
export $(cut -d= -f1 [.env file path])

# Build & run containers
docker compose up -d postgres meilisearch bert mockcas adminer

# Migrate database to v1-alpha version
./scripts/migrate.sh ./migrates/elt-tables recsis-postgres
./scripts/migrate.sh ./migrates/v1-alpha recsis-postgres

# Populate database with data from SIS
docker compose up --build elt

# Init Meilisearch
./scripts/init-meili.sh

# Install templ and wgo
go install github.com/a-h/templ/cmd/templ@v0.2.793
go install github.com/bokwoon95/wgo@latest

# Run webapp
./scripts/watch.sh
```