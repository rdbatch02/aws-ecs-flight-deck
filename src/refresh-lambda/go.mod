module github.com/rdbatch02/ecs-flight-deck/refresh-lambda

go 1.15

require (
	github.com/aws/aws-sdk-go v1.36.6
	github.com/rdbatch02/ecs-flight-deck/domain v0.0.0-00010101000000-000000000000
)

replace github.com/rdbatch02/ecs-flight-deck/domain => ../domain
