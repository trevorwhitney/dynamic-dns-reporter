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

cat > /usr/local/bin/update-dns.sh <<EOF
#!/usr/bin/env bash
set -e
${GOPATH:-$HOME/go}/bin/dynamic-dns-reporter ${accountId} ${apiKey} cerebral
EOF

chmod 700 /usr/local/bin/update-dns.sh
