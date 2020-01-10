# Copyright 2020 Panther Labs Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

AWS_PUBLIC_ACCOUNT_ID=349240696275
AWS_PUBLIC_ECR_REGION=us-west-2
IMAGE_REPO_NAME=panther-ui
PANTHER_VERSION=v0.1

echo Logging in to Amazon ECR...
$(aws ecr get-login --registry-ids $AWS_PUBLIC_ACCOUNT_ID --no-include-email --region $AWS_PUBLIC_ECR_REGION)
echo Build started on `date`
echo Building the Docker image...
docker build -t $IMAGE_REPO_NAME --build-arg PANTHER_VERSION=$PANTHER_VERSION .
docker tag $IMAGE_REPO_NAME $AWS_PUBLIC_ACCOUNT_ID.dkr.ecr.$AWS_PUBLIC_ECR_REGION.amazonaws.com/$IMAGE_REPO_NAME:latest
docker tag $IMAGE_REPO_NAME $AWS_PUBLIC_ACCOUNT_ID.dkr.ecr.$AWS_PUBLIC_ECR_REGION.amazonaws.com/$IMAGE_REPO_NAME:$PANTHER_VERSION
echo Build completed on `date`
echo Pushing the Docker image...
docker push $AWS_PUBLIC_ACCOUNT_ID.dkr.ecr.$AWS_PUBLIC_ECR_REGION.amazonaws.com/$IMAGE_REPO_NAME:latest
docker push $AWS_PUBLIC_ACCOUNT_ID.dkr.ecr.$AWS_PUBLIC_ECR_REGION.amazonaws.com/$IMAGE_REPO_NAME:$PANTHER_VERSION
echo Image uploaded on `date`
