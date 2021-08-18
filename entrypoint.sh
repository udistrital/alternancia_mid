#!/usr/bin/env bash

set -e
set -u
set -o pipefail

# if [ -n "${PARAMETER_STORE:-}" ]; then
#   export ALTERNANCIA_MID__PGUSER="$(aws ssm get-parameter --name /${PARAMETER_STORE}/alternancia_mid/db/username --output text --query Parameter.Value)"
#   export ALTERNANCIA_MID__PGPASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/alternancia_mid/db/password --output text --query Parameter.Value)"
# fi

exec ./main "$@"