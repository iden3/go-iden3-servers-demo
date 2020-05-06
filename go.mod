module github.com/iden3/go-iden3-servers-demo

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.13
	github.com/gin-gonic/gin v1.5.0
	github.com/go-pg/pg/v9 v9.1.2
	github.com/gorilla/mux v1.7.4
	github.com/iden3/go-iden3-core v0.0.8-0.20200428175310-d4796a4acdbe
	github.com/iden3/go-iden3-servers v0.0.2-0.20200505172010-1d41969df64b
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli v1.22.2
	gopkg.in/go-playground/assert.v1 v1.2.1
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/go-playground/validator.v9 v9.29.1
)

// replace github.com/iden3/go-iden3-core => ../go-iden3-core

// replace github.com/iden3/go-iden3-servers => ../go-iden3-servers
