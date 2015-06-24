export DATABASE="user=postgres password=root dbname=embercrest sslmode=disable"

$GOPATH/bin/goose down
$GOPATH/bin/goose up

go test -v ./datastore/pgsql
