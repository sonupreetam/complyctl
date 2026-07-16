#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# Validate EvaluationLog YAML files against the Gemara CUE schema.
#
# Usage:
#   ./tests/schema-validation/validate.sh <scan-output-dir>
#
# The script finds all evaluation-log-*.yaml files in the given directory,
# runs `cue vet` against the Gemara EvaluationLog schema, and reports
# which constraint (if any) failed for each file.
#
# Exit codes:
#   0 — all files pass validation
#   1 — one or more files failed schema validation
#   2 — no evaluation log files found (configuration error)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SCHEMA_DIR="${SCRIPT_DIR}"

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
    output=$(cue vet "${SCHEMA_DIR}/schema.cue" "${file}" 2>&1) || rc=$?

    if [[ ${rc} -eq 0 ]]; then
        PASSED=$((PASSED + 1))
        echo "    PASS"
    else
        FAILED=$((FAILED + 1))
        echo "    FAIL: schema constraint violated"
        echo "    ---"
        echo "${output}" | sed 's/^/    /'
        echo "    ---"
    fi
}

echo "=== EvaluationLog Schema Validation ==="
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
