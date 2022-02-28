#!/usr/bin/env bash
set -e
set -o pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
project_dir="$(cd "$script_dir/.." && pwd)"

credentials=$(op get item dnsimple.com --fields apiToken,accountId)
accountId="$(echo "${credentials}" | jq -r '.accountId')"
apiKey="$(echo "${credentials}" | jq -r '.apiToken')"

pushd "${project_dir}" > /dev/null || exit 1
  go install
popd > /dev/null || exit 1

cat > "${HOME}/.config/systemd/user/dynamicdns.timer" <<EOF
[Unit]
Description=Hourly Dynamic Dns Updater

[Timer]
OnCalendar=daily
AccuracySec=1h
Persistent=true

[Install]
WantedBy=timers.target
EOF

cat > "${HOME}/.config/systemd/user/dynamicdns.service" <<EOF
[Unit]
Description=Dynamic Dns Updater
Wants=dynamicdns.timer

[Service]
Type=oneshot
ExecStart=${GOPATH:-$HOME/go}/bin/dynamic-dns-reporter ${accountId} ${apiKey} ""

[Install]
WantedBy=default.target
EOF

systemctl --user daemon-reload
systemctl --user enable dynamicdns.service
systemctl --user start dynamicdns.service
