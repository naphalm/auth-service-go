#!/bin/sh
set -e

# Name of your Helm release
RELEASE_NAME="auth-service"
# Namespace where the service is deployed
NAMESPACE="default"

echo "Rolling back Helm release: $RELEASE_NAME in namespace $NAMESPACE"

# Rollback to previous revision (1 = last successful deployment)
helm rollback $RELEASE_NAME 1 --namespace $NAMESPACE

echo "Rollback completed."
