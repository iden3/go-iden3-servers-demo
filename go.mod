module github.com/iden3/go-iden3-servers-demo

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.13
	github.com/gin-gonic/gin v1.5.0
	github.com/gorilla/mux v1.7.4
	github.com/iden3/go-circom-prover-verifier v0.0.0-20200522153011-ec6920aa1169
	github.com/iden3/go-iden3-core v0.0.8-0.20200522154712-780c4e804a1a
	github.com/iden3/go-iden3-servers v0.0.2-0.20200525092853-167fda8ae7d9
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli v1.22.2
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/yaml.v2 v2.2.4 // indirect
	xorm.io/xorm v1.0.1
)

// replace github.com/iden3/go-iden3-core => ../go-iden3-core

// replace github.com/iden3/go-iden3-servers => ../go-iden3-servers
