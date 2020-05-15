module github.com/iden3/go-iden3-servers-demo

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.13
	github.com/gin-gonic/gin v1.5.0
	github.com/go-pg/pg/v9 v9.1.2
	github.com/gorilla/mux v1.7.4
	github.com/iden3/go-circom-prover-verifier v0.0.0-20200515100033-bedd64cc7062
	github.com/iden3/go-iden3-core v0.0.8-0.20200515134003-99b3cc33f463
	github.com/iden3/go-iden3-servers v0.0.2-0.20200515143026-71f0da093a96
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli v1.22.2
	gopkg.in/go-playground/assert.v1 v1.2.1
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/go-playground/validator.v9 v9.29.1
	xorm.io/xorm v1.0.1
)

// replace github.com/iden3/go-iden3-core => ../go-iden3-core

// replace github.com/iden3/go-iden3-servers => ../go-iden3-servers
