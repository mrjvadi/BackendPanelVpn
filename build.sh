#!/usr/bin/env bash
set -euo pipefail

REPO_URL="https://github.com/mrjvadi/xrayplugin.git"
TEMP_DIR=$(mktemp -d)
PLUGINS_DIR="./plugins"

echo "📥 Cloning repo..."
git clone --depth 1 "$REPO_URL" "$TEMP_DIR"

mkdir -p "$PLUGINS_DIR"

echo "🔍 Searching for xray/panels/*/client.go..."
find "$TEMP_DIR/xray/panels/" -type f -name "client.go" | while read -r file; do
  panel_dir=$(dirname "$file")
  version=$(basename "$panel_dir")
  panel=$(basename "$(dirname "$panel_dir")")

  safe_panel=${panel// /_}
  so_name="${safe_panel}_${version}.so"
  out_path="${PLUGINS_DIR}/${so_name}"

  echo "🔨 Building plugin for '$panel' v'$version' → $out_path"

  GOOS=linux GOARCH=amd64 \
    go build -buildmode=plugin \
      -o "$out_path" \
      "$panel_dir"
done

echo "🧹 Cleaning up..."
rm -rf "$TEMP_DIR"

echo "✅ All plugins built into $PLUGINS_DIR"

