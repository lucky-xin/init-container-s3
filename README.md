# k8s container 初始化文件下载

## Docker镜像构建
运行脚本
```shell
build-imaga.sh
```
## 环境变量配置
### 环境变量配置-minio 样例：

| Key                     | Description                      | Default Value |
|-------------------------|----------------------------------|---------------|
| `AWS_ACCESS_KEY_ID`     | minio Access Key                 | null          |
| `AWS_SECRET_ACCESS_KEY` | minio Secret Key                 | null          |
| `S3_ENDPOINT`           | 服务器地址：127.0.0.1:9000             | null          |
| `S3_BUCKET`             | 主题 ：piston-algo-model            | null          |
| `S3_KEY`                | 前缀 ： 项目名称，comment-carname/carner | null          |
| `S3_TARGET`             | 目标路径：/home/model                 | null          |
| `S3_TYPE`               | 类型： 目录=dir，文件=file               | dir           |
| `S3_REGION`             | 区域： s3 规范的存储区域                   | us-east-1     |
| `S3_FORCE_PATH_STYLE`   | 路径规则： 兼容模式，阿里云必须为false           | false         |

### 环境变量配置-OSS 样例：
| Key                     | Description                                                                               | Default Value |
|-------------------------|-------------------------------------------------------------------------------------------|---------------|
| `AWS_ACCESS_KEY_ID`     | 阿里云 Access Key                                                                            | null          |
| `AWS_SECRET_ACCESS_KEY` | 阿里云 Secret Key                                                                            | null          |
| `S3_ENDPOINT`           | 服务器地址：<br>内网访问：oss-cn-shenzhen-internal.aliyuncs.com<br>外网访问：oss-cn-shenzhen.aliyuncs.com | null          |
| `S3_BUCKET`             | 存储空间名称 ：piston-algo-model                                                                 | null          |
| `S3_KEY`                | 前缀 ： 项目名称，comment-carname/carner                                                          | null          |
| `S3_TARGET`             | 目标路径：/home/model                                                                          | null          |
| `S3_TYPE`               | 类型： 目录=dir，文件=file                                                                        | dir           |
| `S3_REGION`             | 区域： s3 规范的存储区域                                                                            | us-east-1     |
| `S3_FORCE_PATH_STYLE`   | 路径规则： 兼容模式,阿里云必须为false                                                                    | false         |

### Docker环境运行
```shell
docker run -it -rm \
-e AWS_ACCESS_KEY_ID="" \
-e AWS_SECRET_ACCESS_KEY="" \
-e S3_ENDPOINT="127.0.0.1:9000" \
-e S3_BUCKET='xxx-algo-model' \
-e S3_KEY="/" \
-e S3_TARGET=/home/model \
-e S3_TYPE='' \
-e S3_REGION="admin" \
-v /mnt/d/temp/models:/home/model \
xyz.com/library/container/init-container-s3:1.0.8
```

### k8s环境运行
```yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: echo-test-with-s3-file-init-container
  annotations:
    deployment.kubernetes.io/revision: '1'
spec:
  replicas: 1
  selector:
    matchLabels:
      app: echo-test-with-s3-file-init-container
  template:
    metadata:
      labels:
        app: echo-test-with-s3-file-init-container
    spec:
      initContainers:
        - name: download-s3-file-container
          image: xyz.com/library/container/init-container-s3:1.0.8
          env:
            - name: S3_ENDPOINT
              value: 127.0.0.1:9000
            - name: S3_BUCKET
              value: piston-algo-model
            - name: S3_KEY
              value: rv-valuation
            - name: S3_TARGET
              value: /tmp/rv-valuation
            - name: S3_TYPE
              # dir or file
              value: dir
            - name: S3_REGION
              # default us-east-1
              value: us-east-1
          envFrom:
            - secretRef:
                name: s3-secret
          volumeMounts:
            - name: valuation-model-file
              mountPath: /tmp/rv-valuation
      containers:
        - name: echo-test-with-s3-file-init-container
          image: radial/busyboxplus:curl
          command:
            - /bin/sh
            - '-c'
            - >-
              count=1;while true;do echo log to stdout $count;sleep
              1;count=$(($count+1));done
          resources: { }
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - name: valuation-model-file
              mountPath: /tmp/rv-valuation

      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: { }
      schedulerName: default-scheduler
      volumes:
        - name: valuation-model-file
          emptyDir:
            { }
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 25%
  revisionHistoryLimit: 10
  progressDeadlineSeconds: 600
---
kind: Secret
apiVersion: v1
metadata:
  name: s3-secret
data:
  AWS_ACCESS_KEY_ID: []
  AWS_SECRET_ACCESS_KEY: []
```
