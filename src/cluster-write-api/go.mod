module github.com/rdbatch02/ecs-flight-deck/cluster-write-api

go 1.15

require (
	github.com/aquasecurity/lmdrouter v0.3.0 // indirect
	github.com/aws/aws-lambda-go v1.22.0
	github.com/aws/aws-sdk-go v1.36.20 // indirect
	github.com/rdbatch02/ecs-flight-deck/domain v0.0.0-00010101000000-000000000000
)

replace github.com/rdbatch02/ecs-flight-deck/domain => ../domain
