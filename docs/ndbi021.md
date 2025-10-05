# NDBI021

## Setup

1. Install [Docker Engine](https://docs.docker.com/engine/install/) or [Docker Desktop](https://www.docker.com/products/docker-desktop/)
2. Clone the repository branch
    - `git clone -b ndbi021/main git@github.com:michalhercik/RecSIS.git`
3. Download data
    1. Download data from [OneDrive]()
    2. Unzip
    3. Move `70-load-data.sql` to `RecSIS/init_db`
    4. Move `data.ms.snapshot` to `RecSIS/meili_data/...`
4. Initialize environment and run services

```
scripts/init-env.ps1
docker compose up ...
scripts/init-meili.ps1
```

## Run

```
scripts/watch.ps1
```

## Extend

```
cd recommender
cp algo/template.py algo/[name].py
```

![recommender-schema](recommender-schema.svg)