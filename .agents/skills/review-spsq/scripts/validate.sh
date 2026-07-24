#!/usr/bin/env bash

set -u

if (( $# == 0 )); then
  printf 'usage: %s FILE.spsq [FILE.spsq ...]\n' "$0" >&2
  exit 2
fi

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
repo_root="$(cd -- "${script_dir}/../../../.." && pwd -P)"

validator_kind=""
validator_path=""

if [[ -x "${repo_root}/bin/synapseq" ]]; then
  validator_kind="binary"
  validator_path="${repo_root}/bin/synapseq"
elif command -v synapseq >/dev/null 2>&1; then
  validator_kind="binary"
  validator_path="$(command -v synapseq)"
elif command -v go >/dev/null 2>&1 && [[ -f "${repo_root}/go.mod" ]]; then
  validator_kind="go-run"
else
  printf 'error: no SynapSeq validator is available\n' >&2
  printf 'tried: %s, installed synapseq, and repository go run fallback\n' \
    "${repo_root}/bin/synapseq" >&2
  exit 127
fi

overall_status=0

for input_path in "$@"; do
  printf '=== %s ===\n' "$input_path"

  if [[ ! -f "$input_path" ]]; then
    printf 'command: not run\n'
    printf 'error: file not found or not a regular file\n' >&2
    printf 'status: 2\n'
    overall_status=1
    continue
  fi

  input_dir="$(cd -- "$(dirname -- "$input_path")" && pwd -P)"
  input_abs="${input_dir}/$(basename -- "$input_path")"

  if [[ "$validator_kind" == "binary" ]]; then
    printf 'command: %q -test %q\n' "$validator_path" "$input_abs"
    if "$validator_path" -test "$input_abs"; then
      command_status=0
    else
      command_status=$?
    fi
  else
    printf 'command: (cd %q && go run ./cmd/synapseq -test %q)\n' \
      "$repo_root" "$input_abs"
    if (cd -- "$repo_root" && go run ./cmd/synapseq -test "$input_abs"); then
      command_status=0
    else
      command_status=$?
    fi
  fi

  printf 'status: %d\n' "$command_status"
  if (( command_status != 0 )); then
    overall_status=1
  fi
done

exit "$overall_status"
