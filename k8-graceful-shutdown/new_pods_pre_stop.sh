#!/bin/bash

# current pods
CURRENT_PODS=($(minikube kubectl -- get pods -o jsonpath='{.items[*].metadata.name}'))
echo "${#CURRENT_PODS[@]} pods are running"
    echo "Current pods: ${CURRENT_PODS[@]}"

# rollout
minikube kubectl -- rollout restart deployment/k8s-gs

# watch pods for 20 seconds
for i in {1..20}; do
  sleep 1
  NEW_PODS=()
  PODS=($(minikube kubectl -- get pods -o jsonpath='{.items[*].metadata.name}'))
  for pod in "${PODS[@]}"; do
    if [[ ! " ${CURRENT_PODS[@]} " =~ " ${pod} " ]]; then
      NEW_PODS+=($pod)
    fi
  done

  if [ ${#NEW_PODS[@]} -eq ${#CURRENT_PODS[@]} ] && [ ${#PODS[@]} -eq 4 ]; then
    echo "${#NEW_PODS[@]} new pods was created when the old pods are still running"
    echo "New pods: ${NEW_PODS[@]}"
    break
  fi
done
