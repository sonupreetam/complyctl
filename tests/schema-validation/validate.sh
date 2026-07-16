#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# Validate EvaluationLog YAML files against the Gemara CUE schema.
# Fetches the schema directly from the CUE Central Registry — no local
# module files to maintain.
#
# Usage:
#   ./tests/schema-validation/validate.sh <scan-output-dir> [scan-output-dir ...]
#
# Exit codes:
#   0 — all files pass validation
#   1 — one or more files failed schema validation
#   2 — no evaluation log files found (configuration error)

set -euo pipefail

# Gemara CUE schema module — pinned to v0.23.0 (31f58fb674fd3f3f088533af9c6cc83a2d84e17f).
# Tracks the v0.x line matching go-gemara v0.7.0 used by complyctl.
GEMARA_CUE_MODULE="${GEMARA_CUE_MODULE:-github.com/gemaraproj/gemara@v0.23.0}"

if [[ $# -lt 1 ]]; then
    echo "Usage: $0 <scan-output-dir> [scan-output-dir ...]"
    exit 2
fi

PASSED=0
FAILED=0
TOTAL=0

validate_file() {
    local file="$1"
    TOTAL=$((TOTAL + 1))
    echo "  Validating: ${file}"

    local output rc=0
    output=$(cue vet -c -d '#EvaluationLog' "${GEMARA_CUE_MODULE}" "${file}" 2>&1) || rc=$?

    if [[ ${rc} -eq 0 ]]; then
        PASSED=$((PASSED + 1))
        echo "    PASS"
    else
        FAILED=$((FAILED + 1))
        echo "    FAIL: schema constraint violated"
        echo "    ---"
        echo "    ${output//$'\n'/$'\n'    }"
        echo "    ---"
    fi
}

echo "=== EvaluationLog Schema Validation ==="
echo "  Module: ${GEMARA_CUE_MODULE}"
echo ""

FILES_FOUND=0
for dir in "$@"; do
    if [[ ! -d "${dir}" ]]; then
        echo "WARNING: directory does not exist: ${dir}"
        continue
    fi

    while IFS= read -r -d '' file; do
        FILES_FOUND=$((FILES_FOUND + 1))
        validate_file "${file}"
    done < <(/usr/bin/find "${dir}" -name 'evaluation-log-*.yaml' -print0 2>/dev/null)
done

echo ""

if [[ ${FILES_FOUND} -eq 0 ]]; then
    echo "ERROR: No evaluation-log-*.yaml files found in: $*"
    echo "       Ensure complyctl scan completed successfully before validation."
    exit 2
fi

echo "==============================="
echo "  Total:  ${TOTAL}"
echo "  Passed: ${PASSED}"
echo "  Failed: ${FAILED}"
echo "==============================="

if [[ ${FAILED} -gt 0 ]]; then
    exit 1
fi
