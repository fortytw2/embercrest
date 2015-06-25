export DATABASE="user=postgres password=root dbname=embercrest sslmode=disable"

echo "preparing database for test"

$GOPATH/bin/goose down
$GOPATH/bin/goose up

go test -v ./datastore/pgsql

echo "clearing database after tests"

$GOPATH/bin/goose down
$GOPATH/bin/goose up
