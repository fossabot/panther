/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { RoleNameEnum, SeverityEnum } from 'Generated/schema';
import { BadgeProps } from 'pouncejs';

export const AWS_ACCOUNT_ID_REGEX = new RegExp('\\d{12}');

export const INCLUDE_DIGITS_REGEX = new RegExp('(?=.*[0-9])');

export const INCLUDE_LOWERCASE_REGEX = new RegExp('(?=.*[a-z])');

export const INCLUDE_UPPERCASE_REGEX = new RegExp('(?=.*[A-Z])');

export const INCLUDE_SPECIAL_CHAR_REGEX = new RegExp('(?=.*[!@#\\$%\\^&\\*;:,.<>?/])');

export const DEFAULT_POLICY_FUNCTION =
  'def policy(resource):\n\t# Write your code here.\n\treturn True';

export const DEFAULT_RULE_FUNCTION = 'def rule(event):\n\t# Write your code here.\n\treturn False';

export const RESOURCE_TYPES = [
  'AWS.ACM.Certificate',
  'AWS.CloudFormation.Stack',
  'AWS.CloudTrail',
  'AWS.CloudTrail.Meta',
  'AWS.CloudWatch.LogGroup',
  'AWS.Config.Recorder',
  'AWS.Config.Recorder.Meta',
  'AWS.DynamoDB.Table',
  'AWS.EC2.AMI',
  'AWS.EC2.Instance',
  'AWS.EC2.NetworkACL',
  'AWS.EC2.SecurityGroup',
  'AWS.EC2.Volume',
  'AWS.EC2.VPC',
  'AWS.ELBV2.ApplicationLoadBalancer',
  'AWS.GuardDuty.Detector',
  'AWS.IAM.Group',
  'AWS.IAM.Policy',
  'AWS.IAM.Role',
  'AWS.IAM.RootUser',
  'AWS.IAM.User',
  'AWS.KMS.Key',
  'AWS.Lambda.Function',
  'AWS.PasswordPolicy',
  'AWS.RDS.Instance',
  'AWS.Redshift.Cluster',
  'AWS.S3.Bucket',
  'AWS.WAF.Regional.WebACL',
  'AWS.WAF.WebACL',
] as const;

export const LOG_TYPES = [
  'AWS.AuroraMySQLAudit.Log',
  'AWS.CloudTrail.Log',
  'AWS.S3ServerAccess.Log',
  'AWS.VPCFlow.Log',
  'Osquery.Batch.Log',
  'Osquery.Differential.Log',
  'Osquery.Snapshot.Log',
  'Osquery.Status.Log',
] as const;

export const SEVERITY_COLOR_MAP: { [key in SeverityEnum]: BadgeProps['color'] } = {
  [SeverityEnum.Critical]: 'red' as const,
  [SeverityEnum.High]: 'pink' as const,
  [SeverityEnum.Medium]: 'blue' as const,
  [SeverityEnum.Low]: 'grey' as const,
  [SeverityEnum.Info]: 'neutral' as const,
};

export const PANTHER_SCHEMA_DOCS_LINK = 'https://docs.runpanther.io';

export const DEFAULT_SMALL_PAGE_SIZE = 10;
export const DEFAULT_LARGE_PAGE_SIZE = 25;

// The key under which User-related data will be stored in the storage
export const USER_INFO_STORAGE_KEY = 'panther.user.info';

export const READONLY_ROLES_ARRAY = [RoleNameEnum.ReadOnly];
export const ADMIN_ROLES_ARRAY = [RoleNameEnum.Admin];

export enum INTEGRATION_TYPES {
  AWS_LOGS = 'aws-s3',
  AWS_INFRA = 'aws-scan',
}

export const PANTHER_AUDIT_ROLE = 'panther-compliance-iam';
export const PANTHER_LOG_PROCESSING_ROLE = 'panther-log-processing-role';
export const PANTHER_REAL_TIME = 'panther-cloudwatch-events';
export const PANTHER_REMEDIATION_MASTER_ACCOUNT = 'panther-aws-remediations-master-account';
export const PANTHER_REMEDIATION_SATELLITE_ACCOUNT = 'panther-aws-remediations-satellite-account';
