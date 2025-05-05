#!/bin/bash

# 该版本号应与`deployment`中对应项目的版本号一致
devVer="1.2.0"
notify_url="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=5420e972-c1a7-4a21-b9e4-cda3cb3c42f5"

successMsg() {
  echo -e "\033[32m $1  ✔\033[0m"
}

failMsg() {
  echo -e "\033[31m $1 ✖\033[0m"
}

# 计算$1 到当前时间的差值
cost_time() {
  time_start="$1"
  time_end=$(date "+%Y-%m-%d %H:%M:%S")

  # 将时间字符串转换为时间戳（秒数）
  start_timestamp=$(date -u -j -f "%Y-%m-%d %H:%M:%S" "${time_start}" "+%s" 2>/dev/null)
  end_timestamp=$(date -u -j -f "%Y-%m-%d %H:%M:%S" "${time_end}" "+%s" 2>/dev/null)
  # 检查日期格式是否正确
  if [[ -z "$start_timestamp" || -z "$end_timestamp" ]]; then
      start_timestamp=$(date -d "$time_start" +%s)
      end_timestamp=$(date -d "$time_end" +%s)
  fi

  # 计算时间差值（秒数）
  time_diff=$((end_timestamp - start_timestamp))
  duration=$(echo ${time_diff} | awk '{t=split("60 s 60 m 24 h 999 d",a);for(n=1;n<t;n+=2){if($1==0)break;s=$1%a[n]a[n+1]s;$1=int($1/a[n])}print s}')
  echo "${duration}"
}

# $1 项目, $2 时间 $3 作者, $4 提交信息 , $5成功,失败
noticeDeploySuccess() {
  cost=$(cost_time "$2")
  message=$(cat<< EOF
              {
           "msgtype": "markdown",
          "markdown": {
            "content":"**[$1]自动构建 <font color=\"info\">成功</font>**

构建版本: <font color=\"info\">${devVer}</font>

耗时: ${cost}

作者:  <font color=\"warning\">$3</font>

提交信息: <font color=\"info\">$4</font>"
          }
              }
EOF
)

  curl -X POST -H "Content-Type: application/json" -d "${message}" "${notify_url}"
}

# $1 项目, $2 时间 $3 作者, $4 提交信息
noticeDeployFail() {
  cost=$(cost_time "$2")
  message=$(
    cat <<EOF
              {
           "msgtype": "markdown",
          "markdown": {
            "content":"**[$1]自动构建 <font color=\"warning\">失败</font>**

构建版本: <font color=\"info\">${devVer}</font>

耗时: ${cost}

作者: <font color=\"warning\">$3</font>

提交信息: <font color=\"info\">$4</font>

失败原因: **$5 失败**"
          }
              }
EOF
)

  curl -X POST -H "Content-Type: application/json" -d "${message}" "${notify_url}"
}

# 解析服务名
git_message="$1"
git_author="$2"
git_date=${3%" +0800"}

commit_message="${git_message%%build_*}"

# 构建
docker build -t "registry.cn-hangzhou.aliyuncs.com/feebook/commonapi:${devVer}" .
if [ $? -ne 0 ]; then
  failMsg "构建镜像失败"
  noticeDeployFail   "commonapi" "${git_date}" "${git_author}" "${commit_message}" "构建镜像"
  exit 1
fi

# up
cd /home/app/deployment
git pull
cd /home/app/deployment/docker-compose/test/
docker compose up -d commonapi
if [ $? -ne 0 ]; then
  failMsg "构建镜像失败"
  noticeDeployFail   "commonapi" "${git_date}" "${git_author}" "${commit_message}" "up镜像"
  exit 1
fi

noticeDeploySuccess  "commonapi" "${git_date}" "${git_author}" "${commit_message}"

# 清理镜像
echo -e "清理无效镜像"
docker images | grep "<none>" | awk '{print $3}' | xargs docker rmi
if [ $? -ne 0 ]; then
  failMsg "清理无效镜像"
  exit 0
fi

successMsg "清理无效镜像"
exit 0