#! /bin/bash -e

# Get the list of all .bats files in the current directory
test_files=$(ls ./integration-tests/*.bats)

# Iterate over each test file and run it
for file in $test_files; do
  echo "Running tests in $file"
  docker run -v "${PWD}:/workspace" \
    -e INTEGRATION_TEST_ROOT=/workspace \
    -e GOOS=linux \
    -e GOARCH=$(go env GOARCH) \
    -e TAG_VERSION=$TAG_VERSION \
    -e GITHUB_SHA=$GITHUB_SHA \
    bats/bats:1.11.0 \
    /workspace/integration-tests/$(basename "$file")
done
