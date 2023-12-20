# todolist

A simple Go API to handle todo list.

## Build

`make build`

## Run 

`TASK_API_CONFIG=./config/ ./todolist`

## Environnement
This project need a PostgreSQL database or a MongoDB database.

### PostgreSQL
You can edit the `config.yaml` file or create yours with your connection's data.

Exemple of DSN : 
```
host=<your host> port=<your port> user=<your user> password=<your password> dbname=<your database> sslmode=disable
```

### MongoDB
You can edit the `config.yaml` file or create yours with your connection's data.

Exemple of config : 
```YAML
mongo:
    md1:
        db: <your database>
        host: <your host>
        port: <your port>
        username: <your user>
        password: <your password>
        timeout: 10
```

## Database
### PostgreSQL
The project postgres database is the following :
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
### MongoDB
```JSON
{
    "task_uuid": "<the task UUID>",
    "description": "<the task description>",
    "status": "<the task status>",
    "created_at": "<the task creation time>",
    "last_updated": "<the task last update time>",
}
``````