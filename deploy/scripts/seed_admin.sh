#!/usr/bin/env bash
set -euo pipefail

set -a
source /etc/shieldpanel/shieldpanel.env
set +a
/opt/shieldpanel/bin/shieldpanel-seed
