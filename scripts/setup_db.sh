docker rm -f gochat-pg
docker run --name=gochat-pg -e POSTGRES_PASSWORD='postgres' -p 5432:5432 -d  postgres
sleep 5
docker exec -it gochat-pg psql -U postgres -c 'CREATE DATABASE gochat;'
migrate -database 'postgres://postgres:postgres@localhost:5432/gochat?sslmode=disable' -path schema up