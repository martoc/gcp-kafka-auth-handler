if [ -z "$INTEGRATION_TEST_ROOT" ]; then
  export INTEGRATION_TEST_ROOT=.
fi

BINARY_PATH="$INTEGRATION_TEST_ROOT/target/builds/gcp-kafka-auth-handler-$GOOS-$GOARCH"
