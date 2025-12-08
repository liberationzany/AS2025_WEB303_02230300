module github.com/practical6/tests/integration

go 1.23

require (
	github.com/practical6/menu-service v0.0.0
	github.com/practical6/order-service v0.0.0
	github.com/practical6/proto v0.0.0
	github.com/practical6/user-service v0.0.0
	github.com/stretchr/testify v1.9.0
	google.golang.org/grpc v1.65.0
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.25.12
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240624140628-dc46fd24d27d // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.5.9 // indirect
)

replace github.com/practical6/proto => ../../proto

replace github.com/practical6/user-service => ../../user-service

replace github.com/practical6/menu-service => ../../menu-service

replace github.com/practical6/order-service => ../../order-service
