#!/bin/bash
# Bash strict mode: http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail

# Run linter
docker pull jamescooke/openapi-validator
docker run --user ${UID} --rm -v "$(pwd)/docs:/data" \
-e "NO_UPDATE_NOTIFIER=true" \
jamescooke/openapi-validator \
--errors_only \
--verbose \
openapi.yml

# Clean directory if exists
rm openapi/* || true

# Generate models
docker pull openapitools/openapi-generator-cli
docker run --user ${UID} --rm -v "$(pwd):/local" \
-e "GO_POST_PROCESS_FILE=gofmt -s -w" \
openapitools/openapi-generator-cli generate \
--input-spec /local/docs/openapi.yml \
--generator-name go \
--output /local/openapi \
--additional-properties enumClassPrefix=true

# Remove unused files
find openapi -mindepth 1 -not -iname "model_*.go" -not -name utils.go -delete