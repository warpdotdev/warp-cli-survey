#!/bin/bash
#
# Denver Survey installer
#
# Usage:
#   source <(curl -s https://raw.githubusercontent.com/zachlloyd/denver-survey-client/master/install.sh)

case "$-" in
*i*)	echo This shell is interactive ;;
*)	echo This shell is not interactive ;;
esac

VERSION="0.1.3"

set -e

function run_denver_survey() {
  if [[ "$OSTYPE" == "linux-gnu" ]]; then
      # set -x
      curl -fsSL https://github.com/zachlloyd/denver-survey-client/releases/download/$VERSION/dsurvey.$VERSION.linux.x86_64.tar.gz | tar -xzv dsurvey
      ./dsurvey
  elif [[ "$OSTYPE" == "darwin"* ]]; then
      #set -x
      curl -fsSL https://github.com/zachlloyd/denver-survey-client/releases/download/$VERSION/dsurvey.$VERSION.mac.x86_64.tar.gz | tar -xzv dsurvey
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