# Clean Architecture Golang

Sample project clean architecture implementation in Golang

## Installation

First of all, you need to install Python with pip, GCC and Go. We also need MongoDB and Postgres.  
Say that this project located at `/path/foo/bar/Go/src/github.com/isogram/clean-golang`

- copy `env.dist` to `.env` and set your own config at `.env`
- we may create `virtualenv` of this project for example at `/path/foo/bar/Go/src/github.com/isogram/clean-golang/.venv`
    ```bash
    $ virtualenv .venv
    ```
- activate virtualenv by using
    ```bash
    $ source .venv/bin/activate .
    ```
- go to project root directory then install the requirements & execute the migration
    ```bash
    # install requirements
    $ pip install -r requirements.txt
    # execute migrations
    $ alembic upgrade head
    ```
- compile Go project
    ```bash
    $ make all
    ```
- the compiled binary will be stored in `/path/foo/bar/Go/src/github.com/isogram/clean-golang/bin` . it will generate binary file `api` .
- then you can run service by running
    ```bash
    $ ./bin/api
    ```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Thanks!
