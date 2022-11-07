#!/bin/bash

set -ue

helm lint .
helm package .

# Our helm chart is being hosted on GitHub Pages.
#
# In order to update the index.yaml file in the publish step, we need access to this branch of the repository.
GH_PAGES_DIR="$LD_RELEASE_TEMP_DIR/gh-pages"
git clone -b gh-pages https://github.com/launchdarkly/ld-relay-helm.git "$GH_PAGES_DIR"

# We can use the latest built release, along with the previously published
# index.yaml file to generate an updated index.yaml.
helm repo index . --url https://launchdarkly.github.io/ld-relay-helm --merge "$GH_PAGES_DIR"/index.yaml

cp ld-relay-"$LD_RELEASE_VERSION".tgz "$LD_RELEASE_ARTIFACTS_DIR"

# We can place the .tgz and updated index.yaml files into the release docs
# directory, and releaser will handle updating the gh-pages branch for us.
mv ld-relay-"$LD_RELEASE_VERSION".tgz "$LD_RELEASE_DOCS_DIR"
mv index.yaml "$LD_RELEASE_DOCS_DIR"
