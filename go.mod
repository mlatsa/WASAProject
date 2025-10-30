module github.com/mlatsa/WASAProject

go 1.20

require (
	github.com/ardanlabs/conf v1.5.0
	github.com/femito1/WASA v0.0.0-00010101000000-000000000000
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/websocket v1.5.3
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/sirupsen/logrus v1.9.3
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace github.com/femito1/WASA => ./external/femito1-WASA
