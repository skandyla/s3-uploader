# s3-uploader

Small golang application for testing purposes.


### simple cli testing
`cd build && docker compose up db`  
`export POSTGRES_DSN="postgres://db_user:db_pass@localhost/s3_uploader?sslmode=disable"`  
`go run main.go`  
`curl localhost:8080/liveness`  
`curl localhost:8080/__service/info`  


### golang migrations
While docker-compose psql init script is good enough for local testing, if we want to run SQL migrations at remote hosts, it's better to have external migrations.  

We will use:  
https://github.com/golang-migrate/migrate  

initialize new migration:  
`migrate create -ext sql -dir migrations create_users`  

upply migrations:  
`export PGHOST=localhost`
`migrate -path migrations -database "postgres://db_user:db_pass@$PGHOST/s3_uploader?sslmode=disable" up`  
`migrate -path migrations -database "postgres://db_user:db_pass@$PGHOST/s3_uploader?sslmode=disable" down`  


### Testing against remote docker-compose
`export REMOTE_DOCKER_IP=x.x.x.x`  
`export DOCKER_HOST=ssh://root@${REMOTE_DOCKER_IP}`  

`migrate -path migrations -database "postgres://db_user:db_pass@${REMOTE_DOCKER_IP}/s3_uploader?sslmode=disable" up`  
```
20211223205123/u create_users (656.998494ms)
20211223205308/u refresh_tokens (1.231881717s)
```
