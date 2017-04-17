#!/bin/bash
export PATH=$GOPATH/bin:$PATH

# reflex is used for watch
command -v reflex >/dev/null 2>&1 || {
  echo >&2 "reflex required but not installed.  Please run 'go get -u github.com/cespare/reflex' and try again.";
  exit 1;
}

# https://github.com/julienXX/terminal-notifier can be used for OS X notifications
# when ketchup is rebuilt
command -v "terminal-notifier" >/dev/null 2>&1 && {
  export NOTIFIER_CMD="terminal-notifier -message 'rebuilt app' &&"
}

KETCHUP_DIR=$(git rev-parse --show-toplevel)
cd "$KETCHUP_DIR" || exit

# run in development mode
export KETCHUP_ENV=development

reflex -s \
  -R '^admin/node_modules' \
  -R '^data/' \
  -R '^vendor/' \
  -r '.*.go' \
  -d fancy \
  -- \
  bash -c "rm -f ketchup && go build . && $NOTIFIER_CMD ./ketchup start"
