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
echo "[default]
aws_access_key_id = $INPUT_AWS_ACCESS_KEY_ID
aws_secret_access_key = $INPUT_AWS_SECRET_ACCESS_KEY
region = $INPUT_AWS_DEFAULT_REGION" > ~/.aws/credentials
touch ~/.aws/config
echo "[profile default]
region = $INPUT_AWS_DEFAULT_REGION
output = json " > ~/.aws/config

cat ~/.aws/credentials

echo "setting up mongocli params"
export MCLI_PUBLIC_API_KEY=$INPUT_ATLAS_PUBLIC_KEY_CFN
export MCLI_PRIVATE_API_KEY=$INPUT_ATLAS_PRIVATE_KEY_CFN
export MCLI_ORG_ID=$INPUT_ATLAS_ORG_ID_CFN
export MCLI_OUTPUT=json

echo "setting up env vars"
export ATLAS_PUBLIC_KEY=$INPUT_ATLAS_PUBLIC_KEY_CFN
export ATLAS_PRIVATE_KEY=$INPUT_ATLAS_PRIVATE_KEY_CFN
export ATLAS_ORG_ID=$INPUT_ATLAS_ORG_ID_CFN
echo "env vars:$ATLAS_ORG_ID"

cd cfn-resources
./cfn-testing-helper.sh
cat project/rpdk.log
