# Wallet Core EDA

## How to Run

1. initialize docker containers

```bash
docker compose up -d --build
```

2. create the kakfa topic from the confluence UI

2.1 open the url `http://localhost:9021`

2.2 click on the `Topics` tab

3. run the endpoints in the `api/client.http` file

4. query the database from the sql container

```sh
# 1. get the container id
docker exec -it fullcycle-eda-events-mysql-1 mysql -u root -p
# password: root

# 2. list databases
show databases;

# 3. use the database
use wallet;

# 4. list tables
show tables;

# 5. query the table
select * from accounts;

# 6. add balance to account
update accounts set balance = 100 where id = "a6943d9e-6c5b-4514-b060-793fd0ad53d1";
``` 
