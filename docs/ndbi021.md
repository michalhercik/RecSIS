# NDBI021

## Setup

1. Install [Docker Engine](https://docs.docker.com/engine/install/) or [Docker Desktop](https://www.docker.com/products/docker-desktop/)
2. Clone the repository branch
    - `git clone -b ndbi021/main git@github.com:michalhercik/RecSIS.git`
3. Download data
    1. Download data from [OneDrive](https://cunicz-my.sharepoint.com/:u:/g/personal/81411247_cuni_cz/EUcaiCT79tZBj0qEcC9lRTQBTtv8AN3Iz2gBrIMFZwIt1A?e=WWL5hl)
    2. Unzip it
    3. Move the content inside the root of the repository 
    4. Move `70-load-data.sql` to `RecSIS/init_db`
4. Initialize environment variables
    - Windows
        - `scripts/init-env.ps1 docker.env`
    - Linux
        1. `source docker.env`
        2. `export $(cut -d= -f1 docker.env)`
5. Add Mock CAS to hosts
    - Windows
        - `.\scripts\setup-hosts.ps1`
    - Linux
        - `.\scripts\setup-hosts.sh`
6. Build & run containers
    1. `docker compose build -d webapp recommender postgres meilisearch bert mockcas adminer`
    2. `docker compose up -d webapp recommender postgres meilisearch bert mockcas adminer`
7. Check if it works
    - Go to [localhost:8000](https://localhost:8000)
    - Log in
    - Go to [localhost:8000/recommended](https://localhost:8000/recommended)
    - Select algorithm and see recommendations

## Create new algorithm

1. Go to `recommender` directory
2. Create new algorithm from template
    - Name the file and the class after your cas login to avoid merge conflicts
    - `cp algo/template.py algo/[name].py`
3. Implement `recommend` (required) and `fit` (optional) method
    - You can get inspired by `algo/example.py`
    - Use class `DataRepository` in `data_repository.py` to access the database data
      - accessible via `self.data`
    - Info about the user is passed via `User` class in `user.py`

## Test it

1. Go to [localhost:8000/recommended](https://localhost:8000/recommended)
2. Select your algorithm
3. For debugging use `print` statements and check the logs
    - The logs are accessible in terminal where you are running the recommender
    - If you used docker detached (`-d`) mode, use `docker logs -f recommender`


##  Resources
- meilisearch: [localhost:7700](http://localhost:7700)
    - master key: MASTER_KEY
- adminer: [localhost:8080](http://localhost:8080)
    - db: PostgreSQL
    - server: postgres
    - user: recommender
    - password: recommender
    - db: recsis
- recommender api docs: [localhost:8002/docs](http://localhost:8002/docs)
    
## Schema

![recommender-schema](recommender-schema.svg)