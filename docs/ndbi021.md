# NDBI021

## Setup

1. Install [Docker Engine](https://docs.docker.com/engine/install/) or [Docker Desktop](https://www.docker.com/products/docker-desktop/)
2. Clone the repository branch
    - `git clone -b ndbi021/main git@github.com:michalhercik/RecSIS.git`
3. Download data
    1. Download data from [OneDrive](https://cunicz-my.sharepoint.com/:u:/g/personal/81411247_cuni_cz/EUcaiCT79tZBj0qEcC9lRTQBTtv8AN3Iz2gBrIMFZwIt1A?e=nAj5LF)
    2. Unzip it
    3. Move the content inside the root of the RecSIS repository 
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
3. You can use Test accounts to setup different scenarios
    - You can see recommendations for test accounts from any other account 
    - Test account is any account starting with `test-`
    - You can create new test accounts or modify existing ones
4. For debugging use `print` statements and check the logs
    - The logs are accessible in terminal where you are running the recommender
    - If you used docker detached (`-d`) mode, use `docker logs -f recommender`

## Test accounts

- Informatika se specializací Programování a vývoj software (bachelor, NIPVS19B, 2020)
    - *test-b1w* - First year student
    - *test-b1s* - First year student after winter semester
    - *test-b2w* - Second year student
    - *test-b2s* - Second year student after winter semester
    - *test-b3w* - Third year student
    - *test-b3s* - Third year student after winter semester
- Programování a vývoj software (master, NISD23N, 2023)
    - *test-m1w* - First year master student
    - *test-m1s* - First year master student after winter semester
    - *test-m2w* - Second year master student
    - *test-m2s* - Second year master student after winter semester

## Ideas
- 

##  Resources
- meilisearch: [localhost:7700](http://localhost:7700)
    - master key: MASTER_KEY
- adminer: [localhost:8080](http://localhost:8080)
    - db: PostgreSQL
    - server: postgres
    - user: recommender
    - password: recommender
    - db: recsis
    
## Database/Dataframe Tables

### povinn

| Column     | Description                                |     |
| ---------- | ------------------------------------------ | --- |
| povinn     | course code                                |     |
| pnazev     | course name                                |     |
| panazev    | course name in english                     |     |
| vplatiod   | valid from (academic year)                 |     |
| vplatido   | valid to (academic year)                   |     |
| pfakulta   | faculty code (MFF = 11320)                 |     |
| pgarant    | department code (32-KSI, ...)              |     |
| pvyucovan  | V = taught, N = not taught, Z = cancelled  |     |
| vsemzac    | 1 = winter, 2 = summer, 3 = both semesters |     |
| vsempoc    | number of semesters the course is taught   |     |
| vrozsahpr1 | lecture range in winter                    |     |
| vrozsahcv1 | seminar range in winter                    |     |
| vrozsahpr2 | lecture range in summer                    |     |
| vrozsahcv2 | seminar range in summer                    |     |
| vrvcem     | lecture/seminar range unit                 |     |
| vtyp       | examination type (code)                    |     |
| vebody     | credits                                    |     |
| vucit1     | gurantor 1 (code)                          |     |
| vucit2     | guarantor 2 (code)                         |     |
| vucit3     | guarantor 3 (code)                         |     |

---

### zkous

| Column   | Description                         |     |
| -------- | ----------------------------------- | --- |
| zident   | study ID                            | FK  |
| zskr     | year                                |     |
| zsem     | semester                            |     |
| zpovinn  | course code                         | FK  |
| zmarx    | order number                        |     |
| zroc     | year of study                       |     |
| zbody    | credits                             |     |
| zsplcelk | result (S = passed, N = not passed) |     |

---

### studium

| Column  | Description                       |        |
| ------- | --------------------------------- | ------ |
| soident | student ID                        |        |
| sident  | study ID                          | UNIQUE |
| sfak    | faculty code                      |        |
| sfak2   | secondary faculty code            |        |
| sdruh   | type of study (Bc, Mgr, PhD, ...) |        |
| sobor   | study program code                |        |
| srokp   | year of enrollment                |        |
| sstav   | study status                      |        |
| sroc    | current year of study             |        |
| splan   | degree plan code                  | FK     |

---

### stud_plan

| Column             | Description                                           |     |
| ------------------ | ----------------------------------------------------- | --- |
| code               | course code                                           | FK  |
| interchangeability | course code                                           | FK  |
| bloc_subject_code  | bloc ID                                               |     |
| bloc_type          | A = compulsory, B = compulsory elective, C = elective |     |
| bloc_grade         | recommended year of study                             |     |
| bloc_limit         | number of credits to complete in the block            |     |
| bloc_name_cz       | block name in Czech                                   |     |
| bloc_name_en       | block name in English                                 |     |
| plan_code          | degree plan code                                      |     |
| plan_year          | valid for year                                        |     |

---