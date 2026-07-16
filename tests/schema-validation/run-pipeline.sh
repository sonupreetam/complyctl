#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# Pipeline script for EvaluationLog schema validation CI.
# Runs inside the devcontainer: builds complyctl, installs providers,
# starts the mock OCI registry, and runs get/generate/scan for each
# configured test policy. Scan output is written to a well-known
# directory accessible by the host for subsequent CUE validation.
#
# Expected environment:
#   GITHUB_TOKEN — required for Ampel provider (branch-protection scan)
#   /workspace   — mounted complyctl repository root
#
# Output:
#   /workspace/.ci-scan-output/ — EvaluationLog YAML files

set -euo pipefail

WORKSPACE="/workspace"
SCAN_OUTPUT="/workspace/.ci-scan-output"
REGISTRY_PORT=8765
REGISTRY_URL="http://localhost:${REGISTRY_PORT}"

export PATH="${WORKSPACE}/bin:${GOPATH:-$(go env GOPATH)}/bin:${PATH}"

# ---------------------------------------------------------------------------
# Step 1: Build complyctl and mock-oci-registry
# ---------------------------------------------------------------------------
echo ">>> Building complyctl..."
make -C "${WORKSPACE}" build
echo "    Build complete."

# ---------------------------------------------------------------------------
# Step 2: Install provider dependencies (snappy, ampel, conftest)
# ---------------------------------------------------------------------------
echo ">>> Installing snappy, ampel, and conftest..."
go install github.com/carabiner-dev/snappy@v0.2.4
go install github.com/carabiner-dev/ampel/cmd/ampel@v1.2.1
go install github.com/open-policy-agent/conftest@v0.68.2
echo "    Dependencies installed."

# ---------------------------------------------------------------------------
# Step 3: Clone and build complytime-providers
# ---------------------------------------------------------------------------
echo ">>> Cloning complytime-providers..."
PROVIDERS_TMP="$(mktemp -d)"
trap 'rm -rf "${PROVIDERS_TMP}"' EXIT

git clone --depth 1 \
    https://github.com/complytime/complytime-providers.git \
    "${PROVIDERS_TMP}/complytime-providers"

PROVIDERS_SHA="$(git -C "${PROVIDERS_TMP}/complytime-providers" rev-parse --short HEAD)"
echo "    Cloned at ${PROVIDERS_SHA}"

echo ">>> Building complytime-providers..."
make -C "${PROVIDERS_TMP}/complytime-providers" build

mkdir -p "${HOME}/.complytime/providers"
for provider in ampel openscap opa; do
    binary="complyctl-provider-${provider}"
    src="${PROVIDERS_TMP}/complytime-providers/bin/${binary}"
    if [[ -f "${src}" ]]; then
        cp "${src}" "${HOME}/.complytime/providers/"
        echo "    Installed ${binary}"
    fi
done

# ---------------------------------------------------------------------------
# Step 4: Set up workspace
# ---------------------------------------------------------------------------
echo ">>> Setting up test workspace..."
WORK_DIR="$(mktemp -d)"
TESTDATA_DIR="${WORKSPACE}/tests/cross-repo/testdata"

mkdir -p "${WORK_DIR}/.complytime/ampel/granular-policies"
cp "${TESTDATA_DIR}/complytime.yaml" "${WORK_DIR}/.complytime/"
cp "${TESTDATA_DIR}/granular-policies/block-force-push.json" \
    "${WORK_DIR}/.complytime/ampel/granular-policies/"

cat > "${WORK_DIR}/test-deployment.yaml" << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
        - name: web
          image: nginx:1.27
          securityContext:
            runAsNonRoot: true
          resources:
            limits:
              cpu: "500m"
              memory: "128Mi"
EOF
echo "    Test workspace: ${WORK_DIR}"

# ---------------------------------------------------------------------------
# Step 5: Start mock OCI registry
# ---------------------------------------------------------------------------
echo ">>> Starting mock OCI registry..."
GEMARA_SERVICE_PORT="${REGISTRY_PORT}" "${WORKSPACE}/bin/mock-oci-registry" &
REGISTRY_PID=$!

for _ in $(seq 1 30); do
    if curl -sf "${REGISTRY_URL}/v2/" >/dev/null 2>&1; then
        break
    fi
    sleep 0.5
done
if ! curl -sf "${REGISTRY_URL}/v2/" >/dev/null 2>&1; then
    echo "FATAL: mock registry did not start on ${REGISTRY_URL}"
    exit 1
fi
echo "    Registry ready (PID ${REGISTRY_PID})"

# ---------------------------------------------------------------------------
# Step 6: Run complyctl get
# ---------------------------------------------------------------------------
echo ""
echo ">>> Running complyctl get..."
(cd "${WORK_DIR}" && complyctl get)

# ---------------------------------------------------------------------------
# Step 7: Run generate + scan for each test policy
# ---------------------------------------------------------------------------
POLICIES=("test-ampel-bp" "test-opa-bp")
SCAN_FAILURES=0

for policy_id in "${POLICIES[@]}"; do
    echo ""
    echo ">>> generate --policy-id ${policy_id}..."
    if ! (cd "${WORK_DIR}" && complyctl generate --policy-id "${policy_id}"); then
        echo "    WARNING: generate failed for ${policy_id}, skipping scan."
        continue
    fi

    echo ">>> scan --policy-id ${policy_id}..."
    if ! (cd "${WORK_DIR}" && complyctl scan --policy-id "${policy_id}"); then
        echo "    WARNING: scan returned non-zero for ${policy_id} (may have operational errors)."
        SCAN_FAILURES=$((SCAN_FAILURES + 1))
    fi
done

# ---------------------------------------------------------------------------
# Step 8: Copy scan output to well-known location for host validation
# ---------------------------------------------------------------------------
echo ""
echo ">>> Collecting scan output..."
mkdir -p "${SCAN_OUTPUT}"

SCAN_DIR="${WORK_DIR}/.complytime/scan"
if [[ -d "${SCAN_DIR}" ]]; then
    cp "${SCAN_DIR}"/evaluation-log-*.yaml "${SCAN_OUTPUT}/" 2>/dev/null || true
fi

echo "    Files in ${SCAN_OUTPUT}:"
ls -1 "${SCAN_OUTPUT}/" 2>/dev/null || echo "    (empty)"

# Cleanup
kill "${REGISTRY_PID}" 2>/dev/null || true

EVAL_COUNT=$(/usr/bin/find "${SCAN_OUTPUT}" -name 'evaluation-log-*.yaml' 2>/dev/null | wc -l)
echo ""
echo ">>> Pipeline complete. ${EVAL_COUNT} EvaluationLog file(s) produced."

if [[ ${EVAL_COUNT} -eq 0 ]]; then
    echo "ERROR: No EvaluationLog files were generated. Check scan output above."
    exit 1
fi
