mkdir myapi
Set-Location myapi
go mod init myapi
go get -u github.com/gin-gonic/gin
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
go get -u github.com/golang-jwt/jwt/v4
go get -u github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/gin-contrib/cors
go get -u github.com/lib/pq
