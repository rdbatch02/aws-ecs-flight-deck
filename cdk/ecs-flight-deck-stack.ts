import * as cdk from '@aws-cdk/core';
import { AttributeType, BillingMode, Table } from '@aws-cdk/aws-dynamodb'
import { FlightDeckConfig } from './flight-deck-config';
import { AssetCode, Code, Function, Runtime, Tracing } from '@aws-cdk/aws-lambda';
import { Effect, PolicyStatement, Role, ServicePrincipal } from '@aws-cdk/aws-iam';
import { HttpProxyIntegration, LambdaProxyIntegration } from '@aws-cdk/aws-apigatewayv2-integrations'
import { HttpApi, HttpMethod } from '@aws-cdk/aws-apigatewayv2';
import { Queue } from '@aws-cdk/aws-sqs';
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
import { Bucket, BucketEncryption } from '@aws-cdk/aws-s3';
import { BucketDeployment, Source } from '@aws-cdk/aws-s3-deployment';
import { isMainThread } from 'worker_threads';
import { write } from 'fs';

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
    
    const writeApiLambda = new Function(this, 'write-api-function', {
      functionName: 'EcsFlightDeck-WriteApi',
      runtime: Runtime.GO_1_X,
      memorySize: 512,
      code: AssetCode.fromAsset("src/cluster-write-api/build/writeapi.zip"),
      handler: "WriteApi",
      tracing: Tracing.ACTIVE,
      environment: {
        "REFRESH_CLOCK_LAMBDA_NAME": refreshClockLambda.functionName
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

    const writeApiIntegration = new LambdaProxyIntegration({
      handler: writeApiLambda
    })

    const httpApi = new HttpApi(this, 'EcsFlightDeck-HttpApi')
    httpApi.addRoutes({
      path: '/api/clusters',
      methods: [HttpMethod.GET],
      integration: readApiIntegration
    })

    httpApi.addRoutes({
      path: '/api/clusters',
      methods: [HttpMethod.PUT, HttpMethod.POST],
      integration: writeApiIntegration
    })

    const readEcsPolicy = new PolicyStatement({
      effect: Effect.ALLOW,
      actions: ["ecs:List*", "ecs:Describe*"],
      resources: ["*"]
    })

    const writeEcsPolicy = new PolicyStatement({
      effect: Effect.ALLOW,
      actions: ["ecs:Update*"],
      resources: ["*"]
    })

    refreshLambda.addToRolePolicy(readEcsPolicy)
    readApiLambda.addToRolePolicy(readEcsPolicy)
    writeApiLambda.addToRolePolicy(writeEcsPolicy)

    dbTable.grantReadWriteData(refreshLambda)
    dbTable.grantReadData(readApiLambda)

    // UI

    const bucket = new Bucket(this, 'EcsFlightDeck-s3Bucket', {
      bucketName: 'aws-ecs-flight-deck',
      websiteIndexDocument: 'index.html'
    })

    new BucketDeployment(this, 'DeployUi', {
      sources: [Source.asset('src/ui/build/release')],
      destinationBucket: bucket
    })
    
    const bucketHttpIntegration = new HttpProxyIntegration({
      url: bucket.bucketWebsiteUrl
    })

    httpApi.addRoutes({
      path: '/ui/{proxy+}',
      methods: [HttpMethod.ANY],
      integration: bucketHttpIntegration
    })
  }
}
