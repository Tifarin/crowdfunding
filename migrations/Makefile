up:
	goose postgres "host=localhost port=5432 user=postgres password=password dbname=kubik_test sslmode=disable" up

status:
	goose postgres "user=postgres password=password dbname=kubik_test sslmode=disable" status

reset:
	goose postgres "user=postgres password=password dbname=kubik_test sslmode=disable" reset
	
drop-tables:
	PGPASSWORD=password psql -h localhost -U postgres -d kubik_test -c "DROP SCHEMA public CASCADE;CREATE SCHEMA public;GRANT ALL ON SCHEMA public TO postgres;GRANT ALL ON SCHEMA public TO public;"
