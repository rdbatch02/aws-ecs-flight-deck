module github.com/rdbatch02/ecs-flight-deck/cluster-write-api

go 1.15

require (
	github.com/aws/aws-lambda-go v1.22.0
	github.com/aws/aws-sdk-go v1.36.20
	github.com/awslabs/aws-lambda-go-api-proxy v0.9.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/ugorji/go v1.2.3 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sys v0.0.0-20210113181707-4bcb84eeeb78 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/rdbatch02/ecs-flight-deck/domain => ../domain
