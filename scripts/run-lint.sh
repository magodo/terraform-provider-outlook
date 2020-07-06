#!/usr/bin/env bash

function checkForConditionalRun {
  if [ "$CI" == "true" ];
  then
    echo "Checking if this should be conditionally run.."
    result=$(git diff --name-only origin/master | grep outlook/)
    if [ "$result" = "" ];
    then
      echo "No changes committed to ./outlook - nothing to lint - exiting"
      exit 0
    fi
  fi
}

function runLinters {
  echo "==> Checking source code against linters..."
  golangci-lint run ./...
}

function main {
  checkForConditionalRun
  runLinters
}

main
