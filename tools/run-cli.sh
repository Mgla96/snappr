#!/bin/bash

# This command has already been ran
# ollama run llama3.2 --format json

COMMIT_SHA=$(git rev-parse HEAD)
BRANCH_NAME=qa-${COMMIT_SHA}-$(uuidgen)


LLM_API=ollama
LLM_MODEL=llama3.2
REPOSITORY=snappr
REPOSITORY_OWNER=Mgla96
LLM_ENDPOINT=http://localhost:11434

echo ${BRANCH_NAME}

./snappr-linux-amd64 create \
    --branch ${BRANCH_NAME} \
    --commitSHA ${COMMIT_SHA} \
    --fileRegexPattern ".*\\.go$" \
    --knowledgeSources '' \
    --repository ${REPOSITORY} \
    --repositoryOwner ${REPOSITORY_OWNER} \
    --workflowName createPR
    # --llmAPI ${LLM_API} \
    # --llmEndpoint ${LLM_ENDPOINT} \
    # --llmModel ${LLM_MODEL} \
    

git fetch origin ${BRANCH_NAME}
PR_COMMIT_SHA=$(git rev-parse origin/${BRANCH_NAME})

echo "PR_COMMIT_SHA: ${PR_COMMIT_SHA}"

PR_NUMBER=$(curl -s -H "Authorization: token ${GH_TOKEN}" \
    "https://api.github.com/repos/${REPOSITORY_OWNER}/${REPOSITORY}/pulls?head=${REPOSITORY_OWNER}:${BRANCH_NAME}" \
    | jq -r '.[0].number')

if [ -z "$PR_NUMBER" ]; then
    echo "Pull Request not found!"
    exit 1
fi

echo "PR Number: ${PR_NUMBER}"

./snappr-linux-amd64 review \
    --prNumber ${PR_NUMBER} \
    --commitSHA ${PR_COMMIT_SHA} \
    --fileRegexPattern ".*\\.go$" \
    --repository ${REPOSITORY} \
    --repositoryOwner ${REPOSITORY_OWNER} \
    --workflowName codeReview \
    --fileRegexPattern "ollama.go$"
    # --llmEndpoint ${LLM_ENDPOINT} \
    # --llmAPI ${LLM_API} \
    # --llmModel ${LLM_MODEL}

# then check that comments have been created on the PR
