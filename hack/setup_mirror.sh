#!/usr/bin/env bash
# Setup filesystem mirror for both terraform and tofu. This speeds up the
# integration tests massively.
set -euo pipefail
cd mirror/
terraform providers mirror .
tofu providers mirror .
