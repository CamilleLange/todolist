# todolist

A simple Go API to handle todo list.

## Build

`make build`

## Run 

`TASK_API_CONFIG=./config/ ./todolist`

## Environnement
This project need a database Postgres.

You can edit the `config.yaml` file or create yours with your connection's data.

Exemple of DSN : 
```
host=<your host> port=<your port> user=<your user> password=<your password> dbname=<your database> sslmode=disable
```

## Database
The project database is the following :
```SQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tasks (
    task_uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```