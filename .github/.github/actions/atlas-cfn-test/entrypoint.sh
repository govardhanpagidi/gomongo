#!/bin/bash
# Exit on error. Append "|| true" if you expect an error.
set -o errexit  # same as -e
# Exit on error inside any functions or subshells.
set -o errtrace
# Do not allow use of undefined vars. Use ${VAR:-} to use an undefined VAR
set -o nounset
# Catch if the pipe fucntion fails
set -o pipefail
set -x

echo "#############################################################"
env
#exec "$@"
mkdir -p ~/.aws
touch ~/.aws/credentials
echo "[711489243244_AdministratorAccess]
aws_access_key_id = $INPUT_AWS_ACCESS_KEY_ID
aws_secret_access_key = $INPUT_AWS_SECRET_ACCESS_KEY
region = $INPUT_AWS_DEFAULT_REGION" > ~/.aws/credentials
touch ~/.aws/config
echo "[profile 711489243244_AdministratorAccess]
region = $INPUT_AWS_DEFAULT_REGION
output = json " > ~/.aws/config

export ATLAS_PUBLIC_KEY=$INPUT_ATLAS_PUBLIC_KEY_CFN
export ATLAS_PRIVATE_KEY=$INPUT_ATLAS_PRIVATE_KEY_CFN
export ATLAS_ORG_ID=$INPUT_ATLAS_ORG_ID_CFN

echo "atlas-org-id : $INPUT_ATLAS_ORG_ID_CFN"
echo "atlas_private_key : $INPUT_ATLAS_PRIVATE_KEY_CFN"
echo "atlas_public_key : $INPUT_ATLAS_PUBLIC_KEY_CFN"

cd cfn-resources
BUILD_ONLY=1 CFN_FLAGS="--verbose --set-default --region ${INPUT_AWS_DEFAULT_REGION}" ./cfn-testing-helper.sh
cat project/rpdk.log
