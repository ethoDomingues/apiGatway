module github.com/ethoDomingues/gateway

go 1.23.3

replace github.com/ethoDomingues/c3po => ../c3po

replace github.com/ethoDomingues/braza => ../braza

require (
	github.com/ethoDomingues/braza v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/ethoDomingues/c3po v0.0.0-20240407180005-a2f0a7e9b4ea // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
