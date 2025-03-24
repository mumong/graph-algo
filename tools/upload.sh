#!/bin/bash

# 镜像名称和标签
IMAGE_NAME="algoservice"
IMAGE_TAG="v1.0"
OUTPUT_FILE="${IMAGE_NAME}-${IMAGE_TAG}.tar"

# 工作节点的IP地址
WORKER_NODES=("192.168.0.209" "192.168.0.210")
# 登录凭证
USER="root"
PASS="huhu"

# 1. 检查是否安装了 sshpass
if ! command -v sshpass &> /dev/null
then
    echo "sshpass could not be found. Please install it first."
    exit 1
fi

# 2. 保存 Docker 镜像到本地 tar 文件
echo "Saving Docker image ${IMAGE_NAME}:${IMAGE_TAG} to ${OUTPUT_FILE}..."
docker save -o ${OUTPUT_FILE} ${IMAGE_NAME}:${IMAGE_TAG}

# 3. 将 tar 文件加载到本地主机的 containerd 中
echo "Loading image into local containerd..."
ctr -n k8s.io images import ${OUTPUT_FILE}

# 4. 验证本地是否成功加载
echo "Verifying if the image has been loaded to local containerd..."
crictl images list  -n k8s.io  | grep ${IMAGE_NAME}

# 5. 上传和加载镜像到远程节点
for NODE in "${WORKER_NODES[@]}"; do
    echo "Uploading ${OUTPUT_FILE} to ${NODE}..."
    sshpass -p ${PASS} scp ${OUTPUT_FILE} ${USER}@${NODE}:/tmp/

    echo "Loading image on remote node ${NODE}..."
    sshpass -p ${PASS} ssh ${USER}@${NODE} "ctr -n k8s.io images import /tmp/${OUTPUT_FILE} && rm /tmp/${OUTPUT_FILE}"

    echo "Verifying if the image has been loaded on ${NODE}..."
    sshpass -p ${PASS} ssh ${USER}@${NODE} "crictl images --namespace=k8s.io | grep ${IMAGE_NAME}"

    echo "Image successfully loaded on ${NODE}."
done

echo "Process completed."

# 删除本地 tar 文件
if [ -f "${OUTPUT_FILE}" ]; then
    echo "Deleting local tar file: ${OUTPUT_FILE}"
    rm -f "${OUTPUT_FILE}"
else
    echo "No local tar file found to delete."
fi

echo "Process completed."
