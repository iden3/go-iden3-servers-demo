module github.com/iden3/go-iden3-servers-demo

go 1.13

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/go-pg/pg/v9 v9.1.2
	github.com/iden3/go-iden3-core v0.0.7-0.20200219145505-d0591eb49ae7
	github.com/iden3/go-iden3-servers v0.0.0-20200219150204-1856e87ce0d8
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/go-playground/validator.v9 v9.29.1
)

replace github.com/iden3/go-iden3-core => ../go-iden3-core

replace github.com/iden3/go-iden3-servers => ../go-iden3-servers
