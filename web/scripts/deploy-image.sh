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
