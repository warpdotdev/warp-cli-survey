#!/bin/bash
#
# Denver Survey installer
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/zachlloyd/denver-survey-client/master/install.sh| bash

VERSION="0.1.0"

set -e

function run_denver_survey() {
  if [[ "$OSTYPE" == "linux-gnu" ]]; then
      set -x
      curl -fsSL https://github.com/zachlloyd/denver-survey-client/releases/download/v$VERSION/dsurvey.$VERSION.linux.x86_64.tar.gz | tar -xzv dsurvey
      ./dsurvey
  elif [[ "$OSTYPE" == "darwin"* ]]; then
      set -x
      curl -fsSL https://github.com/zachlloyd/denver-survey-client/releases/download/v$VERSION/dsurvey.$VERSION.mac.x86_64.tar.gz | tar -xzv dsurvey
      ./dsurvey
  else
      set +x
      echo "The Denver Survey installer does not currently work for your platform: $OS"
      echo "Please let me know at zachlloyd@gmail.com"
      exit 1
  fi

  set +x
}

run_denver_survey


