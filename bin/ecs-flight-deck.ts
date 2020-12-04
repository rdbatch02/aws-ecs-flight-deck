#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { EcsFlightDeckStack } from '../lib/ecs-flight-deck-stack';

const app = new cdk.App();
new EcsFlightDeckStack(app, 'EcsFlightDeckStack');
