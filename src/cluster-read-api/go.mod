module github.com/rdbatch02/ecs-flight-deck/cluster-read-api

go 1.15

require (
	github.com/aws/aws-lambda-go v1.21.0 // indirect
	github.com/aws/aws-sdk-go v1.36.15 // indirect
	github.com/aws/aws-xray-sdk-go v1.1.0 // indirect
	github.com/rdbatch02/ecs-flight-deck/domain v0.0.0-00010101000000-000000000000
)

replace github.com/rdbatch02/ecs-flight-deck/domain => ../domain