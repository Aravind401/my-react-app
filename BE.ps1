mkdir myapi
cd myapi
go mod init myapi
go get -u github.com/gin-gonic/gin
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
go get github.com/golang-jwt/jwt/v4
go install github.com/swaggo/swag/cmd/swag@latest
