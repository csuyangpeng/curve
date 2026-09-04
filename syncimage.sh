#!/bin/bash

# 设置变量
HUAWEI_REGISTRY="swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io"
ALI_REGISTRY="registry.cn-hangzhou.aliyuncs.com/10_18_1_2_5000"
IMAGE="redis"
VERSION="8.8"

# 开启实验性功能
export DOCKER_CLI_EXPERIMENTAL=enabled

echo "=== 1. 拉取华为云镜像 ==="
docker pull ${HUAWEI_REGISTRY}/${IMAGE}:${VERSION}
docker pull ${HUAWEI_REGISTRY}/${IMAGE}:${VERSION}-linuxarm64

echo "=== 2. 打标签 ==="
docker tag ${HUAWEI_REGISTRY}/${IMAGE}:${VERSION} ${ALI_REGISTRY}/${IMAGE}:${VERSION}-amd64
docker tag ${HUAWEI_REGISTRY}/${IMAGE}:${VERSION}-linuxarm64 ${ALI_REGISTRY}/${IMAGE}:${VERSION}-arm64

echo "=== 3. 推送镜像到阿里云 ==="
docker push ${ALI_REGISTRY}/${IMAGE}:${VERSION}-amd64
docker push ${ALI_REGISTRY}/${IMAGE}:${VERSION}-arm64

echo "=== 4. 创建并推送 manifest ==="
docker manifest create --amend ${ALI_REGISTRY}/${IMAGE}:${VERSION} \
  ${ALI_REGISTRY}/${IMAGE}:${VERSION}-amd64 \
  ${ALI_REGISTRY}/${IMAGE}:${VERSION}-arm64

docker manifest annotate ${ALI_REGISTRY}/${IMAGE}:${VERSION} \
  ${ALI_REGISTRY}/${IMAGE}:${VERSION}-amd64 --arch amd64 --os linux

docker manifest annotate ${ALI_REGISTRY}/${IMAGE}:${VERSION} \
  ${ALI_REGISTRY}/${IMAGE}:${VERSION}-arm64 --arch arm64 --os linux

docker manifest push --purge ${ALI_REGISTRY}/${IMAGE}:${VERSION}

echo "=== 5. 验证 ==="
docker manifest inspect ${ALI_REGISTRY}/${IMAGE}:${VERSION}

echo "=== 完成！==="