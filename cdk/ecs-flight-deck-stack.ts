import * as cdk from '@aws-cdk/core';
import { AttributeType, BillingMode, Table } from '@aws-cdk/aws-dynamodb'
import { FlightDeckConfig } from './flight-deck-config';

export class EcsFlightDeckStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: FlightDeckConfig) {
    super(scope, id, props);

    const dbTable = new Table(this, 'table', {
      tableName: 'EcsFlightDeck',
      billingMode: BillingMode.PROVISIONED,
      readCapacity: 5,
      writeCapacity: 5,
      partitionKey: {
        name: 'clusterId',
        type: AttributeType.STRING
      }
    })
  }
}
