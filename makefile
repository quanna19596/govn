importdb:
	docker exec -i postgres-db psql -U root -d govn < ./backupdb-govn.sql
exportdb:
	docker exec -i postgres-db pg_dump -U root -d govn > ./backupdb-govn.sql
server:
	go run .