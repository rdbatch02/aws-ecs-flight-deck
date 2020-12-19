import * as cdk from '@aws-cdk/core';
import { AttributeType, BillingMode, Table } from '@aws-cdk/aws-dynamodb'
import { FlightDeckConfig } from './flight-deck-config';
import { AssetCode, Function, Runtime } from '@aws-cdk/aws-lambda';
import { Effect, PolicyStatement } from '@aws-cdk/aws-iam';

export class EcsFlightDeckStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: FlightDeckConfig) {
    super(scope, id, props);

    const dbTable = new Table(this, 'table', {
      tableName: 'EcsFlightDeck-ClusterServices',
      billingMode: BillingMode.PROVISIONED,
      readCapacity: 5,
      writeCapacity: 5,
      partitionKey: {
        name: 'ClusterArn',
        type: AttributeType.STRING
      }
    })

    const refreshLambda = new Function(this, 'refresh-function', {
      functionName: 'EcsFlightDeck-Refresh',
      runtime: Runtime.GO_1_X,
      memorySize: 512,
      code: AssetCode.fromAsset("src/refresh-lambda/refresh.zip"),
      handler: "refresh",
      environment: {
        "FLIGHT_DECK_TABLE_NAME": dbTable.tableName
      }
    })

    refreshLambda.addToRolePolicy(new PolicyStatement({
      effect: Effect.ALLOW,
      actions: ["ecs:List*", "ecs:Describe*"],
      resources: ["*"]
    }))

    dbTable.grantReadWriteData(refreshLambda)
  }
}
