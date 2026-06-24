#!/bin/bash
# set -x
REGISTRY_URL="http://10.18.1.2:5000"
REPO_NAME="react-admin"

# 删除第二位版本号 > 10 的标签，如 2.66、3.22、2.521；非数字版本号跳过
should_delete() {
  local tag="$1"
  if [[ "$tag" =~ ^[0-9]+\.([0-9]+) ]]; then
    local second="${BASH_REMATCH[1]}"
    (( second > 10 ))
    return $?
  fi
  return 1
}

TAGS=$(curl -s "${REGISTRY_URL}/v2/${REPO_NAME}/tags/list" | jq -r '.tags[]?')

if [ -z "$TAGS" ]; then
  echo "未找到任何标签: ${REPO_NAME}"
  exit 1
fi

deleted=0
skipped=0

for TAG in $TAGS; do
  if ! should_delete "$TAG"; then
    echo "跳过: ${REPO_NAME}:${TAG}"
    skipped=$((skipped + 1))
    continue
  fi

  echo "正在删除: ${REPO_NAME}:${TAG}"
  DIGEST=$(curl -s -D - -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
    "${REGISTRY_URL}/v2/${REPO_NAME}/manifests/${TAG}" \
    | grep -i "Docker-Content-Digest" | awk '{print $2}' | tr -d '\r')

  if [ -n "$DIGEST" ]; then
    echo "  Digest: ${DIGEST}"
    curl -s -X DELETE "${REGISTRY_URL}/v2/${REPO_NAME}/manifests/${DIGEST}"
    deleted=$((deleted + 1))
  else
    echo "  获取 Digest 失败，跳过。"
  fi
done

echo "完成: 删除 ${deleted} 个, 跳过 ${skipped} 个。"
echo "请手动执行垃圾回收: docker exec registry registry garbage-collect /etc/docker/registry/config.yml"
