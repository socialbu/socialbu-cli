#!/usr/bin/env sh

set -eu

repository="usamaejaz/socialbu-cli"
version="${SOCIALBU_VERSION:-latest}"
install_dir="${SOCIALBU_INSTALL_DIR:-}"
github_token="${GITHUB_TOKEN:-}"

case "$(uname -s)" in
  Linux) os="linux" ;;
  Darwin) os="darwin" ;;
  *)
    echo "Unsupported operating system: $(uname -s)" >&2
    exit 1
    ;;
esac

case "$(uname -m)" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "Unsupported architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

if [ -z "$install_dir" ]; then
  if [ -d /usr/local/bin ] && [ -w /usr/local/bin ]; then
    install_dir="/usr/local/bin"
  else
    install_dir="${HOME:?HOME is not set}/.local/bin"
  fi
fi

asset="socialbu_${os}_${arch}"
if [ "$version" = "latest" ]; then
  download_base="https://github.com/${repository}/releases/latest/download"
else
  case "$version" in
    v*) tag="$version" ;;
    *) tag="v${version}" ;;
  esac
  download_base="https://github.com/${repository}/releases/download/${tag}"
fi

download_public() {
  source_url="$1"
  destination="$2"
  if command -v curl >/dev/null 2>&1; then
    curl --fail --silent --show-error --location "$source_url" --output "$destination"
  elif command -v wget >/dev/null 2>&1; then
    wget --quiet "$source_url" --output-document="$destination"
  else
    echo "curl or wget is required" >&2
    exit 1
  fi
}

download_release_asset() {
  asset_name="$1"
  destination="$2"
  if [ -z "$github_token" ]; then
    download_public "${download_base}/${asset_name}" "$destination"
    return
  fi
  if ! command -v gh >/dev/null 2>&1; then
    echo "gh is required to install from a private GitHub repository" >&2
    exit 1
  fi

  if [ "$version" = "latest" ]; then
    release_endpoint="repos/${repository}/releases/latest"
  else
    release_endpoint="repos/${repository}/releases/tags/${tag}"
  fi
  asset_id="$(GH_TOKEN="$github_token" gh api "$release_endpoint" \
    --jq ".assets[] | select(.name == \"${asset_name}\") | .id")"
  if [ -z "$asset_id" ]; then
    echo "Release asset not found: ${asset_name}" >&2
    exit 1
  fi
  GH_TOKEN="$github_token" gh api \
    -H "Accept: application/octet-stream" \
    "repos/${repository}/releases/assets/${asset_id}" > "$destination"
}

temp_dir="$(mktemp -d)"
cleanup() {
  rm -rf "$temp_dir"
}
trap cleanup EXIT HUP INT TERM

binary_path="${temp_dir}/${asset}"
checksums_path="${temp_dir}/checksums.txt"
download_release_asset "$asset" "$binary_path"
download_release_asset "checksums.txt" "$checksums_path"

expected_checksum="$(awk -v asset="$asset" '$2 == asset { print $1 }' "$checksums_path")"
if [ -z "$expected_checksum" ]; then
  echo "No checksum found for ${asset}" >&2
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  actual_checksum="$(sha256sum "$binary_path" | awk '{ print $1 }')"
elif command -v shasum >/dev/null 2>&1; then
  actual_checksum="$(shasum -a 256 "$binary_path" | awk '{ print $1 }')"
else
  echo "sha256sum or shasum is required" >&2
  exit 1
fi

if [ "$actual_checksum" != "$expected_checksum" ]; then
  echo "Checksum verification failed for ${asset}" >&2
  exit 1
fi

mkdir -p "$install_dir"
install -m 0755 "$binary_path" "${install_dir}/socialbu"

echo "Installed socialbu to ${install_dir}/socialbu"
case ":${PATH}:" in
  *:"${install_dir}":*) ;;
  *) echo "Add ${install_dir} to PATH before running socialbu." ;;
esac
