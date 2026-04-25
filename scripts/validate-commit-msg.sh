#!/usr/bin/env sh
set -eu

commit_msg_file="${1:-}"

if [ -z "$commit_msg_file" ] || [ ! -f "$commit_msg_file" ]; then
  echo "Usage: scripts/validate-commit-msg.sh <commit-msg-file>" >&2
  exit 2
fi

subject="$(
  sed '/^[[:space:]]*#/d;/^[[:space:]]*$/d;q' "$commit_msg_file"
)"

case "$subject" in
  Merge\ *|Revert\ *|fixup!\ *|squash!\ *)
    exit 0
    ;;
esac

pattern='^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9._/-]+\))?!?: .+'

if printf '%s\n' "$subject" | grep -Eq "$pattern"; then
  exit 0
fi

cat >&2 <<'EOF'
Invalid commit message.

Use Conventional Commits:
  feat: add scaffold prompt suggestions
  fix(api): handle missing plugin
  docs(readme): update local setup
  chore!: drop deprecated route

Allowed types:
  build, chore, ci, docs, feat, fix, perf, refactor, revert, style, test
EOF

exit 1
