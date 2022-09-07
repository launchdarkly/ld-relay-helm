# shellcheck shell=bash

function configure_git() {
    # Configure GitHub access token for pushing to client repositories
    # (We're not doing this in prepare.sh because we want to make sure there's no way the
    # build/test scripts can accidentally push to GitHub)
    echo >>~/.netrc "machine github.com login LaunchDarklyReleaseBot password $(cat "$LD_RELEASE_SECRETS_DIR"/github_token)"

    git config --global user.name LaunchDarklyReleaseBot
    git config --global user.email launchdarklyreleasebot@launchdarkly.com
}

function update_gh_pages() {
    local publish=${1:0}
    local GH_PAGES_DIR="$LD_RELEASE_TEMP_DIR/gh-pages"

    git clone -b gh-pages https://github.com/launchdarkly/ld-relay-helm.git "$GH_PAGES_DIR"
    helm repo index . --url https://launchdarkly.github.io/ld-relay-helm --merge "$GH_PAGES_DIR"/index.yaml

    mv ld-relay-"$LD_RELEASE_VERSION".tgz "$GH_PAGES_DIR"
    mv index.yaml "$GH_PAGES_DIR"/index.yaml

    cd "$GH_PAGES_DIR"
    git add .
    git commit -m "Update chart repository with version $LD_RELEASE_VERSION"

    if [ "$publish" -eq 1 ]; then
        git push origin gh-pages
    fi

    cd -
}

function publish() {
    local publish=${1:0}

    cp ld-relay-"$LD_RELEASE_VERSION".tgz "$LD_RELEASE_ARTIFACTS_DIR"

    configure_git
    update_gh_pages "$publish"
}
