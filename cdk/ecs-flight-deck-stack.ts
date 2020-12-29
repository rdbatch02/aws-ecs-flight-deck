import * as cdk from '@aws-cdk/core';
import { AttributeType, BillingMode, Table } from '@aws-cdk/aws-dynamodb'
import { FlightDeckConfig } from './flight-deck-config';
import { AssetCode, Code, Function, Runtime, Tracing } from '@aws-cdk/aws-lambda';
import { Effect, PolicyStatement } from '@aws-cdk/aws-iam';
import { LambdaProxyIntegration } from '@aws-cdk/aws-apigatewayv2-integrations'
import { HttpApi, HttpMethod } from '@aws-cdk/aws-apigatewayv2';
import { Queue } from '@aws-cdk/aws-sqs';
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';

export class EcsFlightDeckStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: FlightDeckConfig) {
    super(scope, id, props);

    const DELAY_SECONDS = 60

    const dbTable = new Table(this, 'table', {
      tableName: 'EcsFlightDeck-Services',
      billingMode: BillingMode.PROVISIONED,
      readCapacity: 5,
      writeCapacity: 5,
      partitionKey: {
        name: 'ClusterArn',
        type: AttributeType.STRING
      },
      sortKey: {
        name: 'ServiceArn',
        type: AttributeType.STRING
      }
    })

    const refreshQueue = new Queue(this, 'refresh-queue', {
      queueName: 'EcsFlightDeck-refreshQueue'
    })    

    const refreshLambda = new Function(this, 'refresh-function', {
      functionName: 'EcsFlightDeck-Refresh',
      runtime: Runtime.GO_1_X,
      memorySize: 512,
      code: AssetCode.fromAsset("src/refresh-lambda/build/refresh.zip"),
      handler: "refresh",
      environment: {
        "FLIGHT_DECK_TABLE_NAME": dbTable.tableName
      }
    })

    const readApiLambda = new Function(this, 'read-api-function', {
      functionName: 'EcsFlightDeck-ReadApi',
      runtime: Runtime.GO_1_X,
      memorySize: 512,
      code: AssetCode.fromAsset("src/cluster-read-api/build/readapi.zip"),
      handler: "ReadApi",
      tracing: Tracing.ACTIVE,
      environment: {
        "FLIGHT_DECK_TABLE_NAME": dbTable.tableName,
        "REFRESH_QUEUE_URL": refreshQueue.queueUrl,
        "DELAY_SECONDS": DELAY_SECONDS.toString()
      }
    })

    const refreshClockLambda = new Function(this, 'refresh-clock-function', {
      functionName: "EcsFlightDeck-RefreshClock",
      runtime: Runtime.GO_1_X,
      memorySize: 256,
      code: AssetCode.fromAsset("src/refresh-clock/build/refreshclock.zip"),
      handler: "RefreshClock",
      environment: {
        "REFRESH_LAMBDA_NAME": refreshLambda.functionName,
        "REFRESH_QUEUE_URL": refreshQueue.queueUrl,
        "DELAY_SECONDS": (DELAY_SECONDS * 10).toString()
      }
    })

    refreshQueue.grantConsumeMessages(refreshClockLambda)
    refreshQueue.grantSendMessages(refreshClockLambda)
    refreshQueue.grantPurge(refreshClockLambda)

    refreshClockLambda.addEventSource(new SqsEventSource(refreshQueue))

    refreshQueue.grantSendMessages(readApiLambda)

    const readApiIntegration = new LambdaProxyIntegration({
      handler: readApiLambda
    })

    const httpApi = new HttpApi(this, 'EcsFlightDeck-HttpApi')
    httpApi.addRoutes({
      path: '/clusters',
      methods: [HttpMethod.GET],
      integration: readApiIntegration
    })

    const readEcsPolicy = new PolicyStatement({
      effect: Effect.ALLOW,
      actions: ["ecs:List*", "ecs:Describe*"],
      resources: ["*"]
    })

    refreshLambda.addToRolePolicy(readEcsPolicy)
    readApiLambda.addToRolePolicy(readEcsPolicy)

    dbTable.grantReadWriteData(refreshLambda)
    dbTable.grantReadData(readApiLambda)
  }
}
