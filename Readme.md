# OpenFGA based Authz Engine
## What is OpenFGA?
> [OpenFGA](https://openfga.dev/) is an open-source authorization solution that allows developers to build granular access control using an easy-to-read modeling language and friendly APIs.

This PoC hosts a webapp along that connects to an OpenFGA server and adds Authz Models, along with tuple.
It later exposes an endpoint that performs a `check`. The services in this setup can be deployed using `docker-compose`

## How to run the POC
Use `Makefile` rules to get started
1. Run `make run` to kickstart the POC
2. Run `make stop` to bring the service down
3. Run `make clean` to cleanup all containers and images

Use the [dashboard](http://localhost:3000/playground) to view the model, tuples and assertions
![image](dashboard.png)

## Steps to setup
### Authz Model Configuration DSL
OpenFGA Authz Model can be programmed using their custom DSL.
The details can be found [here](https://openfga.dev/docs/configuration-language)
In this repo, you can take a look at `authz-model.fga` to see how we modeled an Impersonation Scenario.
Unfortunately, this DSL would need to be transformed to JSON for the OpenFGA APIs to interact with.
This would can be done using their CLI's [transform](https://github.com/openfga/cli?tab=readme-ov-file#transform-an-authorization-model) option
or using a [VSCode plugin](https://marketplace.visualstudio.com/items?itemName=openfga.openfga-vscode). The plugin does provide auto-completion as well.
The json of `authz-mode.fga` can be found as `authz-model.json` in this same repo.

> #### Note
> The DSL is lot more terse than the JSON.

### Store Setup
Along with model configuration, we would need Stores to be setup as well.
> A store is an OpenFGA entity used to organize authorization check data.
This can be done using CLI and API. In the code, we set this up using the latter.
The webserver setups the model

### Adding Tuples and Assertions to the Store
This also can be done using CLI and APIs. In this PoC, we chose to use the latter. It also adds Assertions to verify the relationships via the [dashboard](http://localhost:3000/playground).
One can view even the Decision Tree using this dashboard for every assertion execution ![image](decision-tree.png). The loading of the tuples can be found in `tuple-loader.go`

### cURL
Finally, we can use `check` to verify that
- an impersonation relationship exists between the imersponator and impersonation
- And, if the impersonator has certain CRUD perms when dealing with a capability.

`/check` endpoint can be used like this:
```bash
curl --location 'localhost:8888/check' --header 'Content-Type: application/json' --data '{
    "user_id": "homer",
    "impersonator_id": "beth",
    "relation": "can_read",
    "capability_id": "claims"
}'
```

When set, the impersonation relation is valid only for `1m`.
We enforce this using CEL-based [Condition](https://openfga.dev/docs/modeling/conditions) that is found in the Authz Model Config.
The condition used to check expiry in this POC is:
```yaml
condition check_expired(current_time: timestamp, grant_time: timestamp, grant_duration: duration) {
  current_time < grant_time + grant_duration
}
```
We finally pass the `current_time` in OpenFGA's `check` API using `Context` map.

### Prometheus Metrics
Prometheus Metrics for the OpenFGA container can be scraped using this [link](http://localhost:2112/metrics)

### Logs
Logs can be viewed by looking at the Openfga container using this command from host:

```bash
docker-compose logs -f openfga 
```

![image](./logs.png)

### Postgres
All the data (model, tuples, assertions) are stored in Postgres. 
The can be accessed after ssh'ing into the postgres container via :
```bash
root@a7c6f8c180b2:/# psql -U postgres
psql (14.13 (Debian 14.13-1.pgdg120+1))
Type "help" for help.

postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
(3 rows)

postgres=# \c postgres
You are now connected to database "postgres" as user "postgres".
postgres=# \dt
                List of relations
 Schema |        Name         | Type  |  Owner   
--------+---------------------+-------+----------
 public | assertion           | table | postgres
 public | authorization_model | table | postgres
 public | changelog           | table | postgres
 public | goose_db_version    | table | postgres
 public | store               | table | postgres
 public | tuple               | table | postgres
(6 rows)

postgres=# select * from tuple limit 5;

           store            | object_type |   object_id    | relation |            _user            | user_type |            ulid            |         inserted_at          | condition_name | condition_context 
----------------------------+-------------+----------------+----------+-----------------------------+-----------+----------------------------+------------------------------+----------------+-------------------
 01J58MKZR32NMREB1WKQ7BVRNN | group       | claims_reader  | member   | user:beth                   | user      | 01J58MKZRFVVN815ZN1GYWS4AS | 2024-08-14 14:28:09.61584+00 |                | 
 01J58MKZR32NMREB1WKQ7BVRNN | capability  | claims         | reader   | group:claims_reader#member  | userset   | 01J58MKZRFVVN815ZN1HWBTQQ2 | 2024-08-14 14:28:09.61584+00 |                | 
 01J58MKZR32NMREB1WKQ7BVRNN | group       | wallet_deleter | member   | user:jerry                  | user      | 01J58MKZRFVVN815ZN1K7CTXDQ | 2024-08-14 14:28:09.61584+00 |                | 
 01J58MKZR32NMREB1WKQ7BVRNN | capability  | wallet         | deleter  | group:wallet_deleter#member | userset   | 01J58MKZRFVVN815ZN1Q50VF2G | 2024-08-14 14:28:09.61584+00 |                | 
 01J58MKZR32NMREB1WKQ7BVRNN | group       | wallet_reader  | member   | user:jerry                  | user      | 01J58MKZRFVVN815ZN1S0KFXC5 | 2024-08-14 14:28:09.61584+00 |                | 
(5 rows)
postgres=# 
```



## TODO
- ~~Add the Impersonator tuple (with CEL condition)programmatically~~ Done
- ~~Expose a Postman collection to do so.~~ Done
- ~~Create a POST to accept impersonation tuples to create~~ Done
- ~~Create a POST /check that accepts json to check if an impersonation is valid and has the correct permissions to perform a CRUD.~~ Done
- ~~See if we can simulate calling into external services similar to http_send in OPA~~ Done
- Make the impersonator relation expiry configurable via Docker-compose.