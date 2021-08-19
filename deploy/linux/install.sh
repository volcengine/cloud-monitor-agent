#!/bin/bash
#火山引擎云监控插件安装脚本
#set -x

echo "installing cloud-monitor-agent"

export LANG=en_US.UTF-8
export LANGUAGE=en_US:

WD_HOME="/usr/local/cloud-monitor-agent"
TAR_FILE=cloud-monitor-agent_linux_${ARCH}_${VERSION}.tar.gz

DEST_TAR_FILE=${WD_HOME}/${TAR_FILE}
DEST_BIN_FILE=${WD_HOME}/cloud-monitor-agent
DEST_MONITOR_CONF=${WD_HOME}/conf_monitor.json
DEST_SERVICE_FILE=/lib/systemd/system/cloud-monitor-agent.service

set_default_region() {
  REGION_ID="LF-BOE"
}

set_default_version() {
  VERSION="v1.0.0"
}

if [[ -z $REGION_ID ]]; then
  set_default_region
fi

if [[ -z $VERSION ]]; then
  set_default_version
fi

if [[ $(uname -m) == "x86_64" ]]; then
  ARCH="amd64"
else
  echo "Unsupported Arch: $(uname -m)"
  exit 1
fi

case $(uname -s) in
Linux)
  WD_OS="linux"
  ;;
*)
  echo "Unsupported OS: $(uname -s)"
  exit 1
  ;;
esac

case ${REGION_ID} in
LF-BOE)
  DOWNLOAD_PATH=http://x/cloud-monitor-agent/${TAR_FILE}
  REGION_CODE=cn-north-4
  ;;
LQ-BOE)
  DOWNLOAD_PATH=http://x/cloud-monitor-agent/${TAR_FILE}
  REGION_CODE=cn-north-3
  ;;
HL-BOE)
  DOWNLOAD_PATH=http://x/cloud-monitor-agent/${TAR_FILE}
  REGION_CODE=cn-north-2
  ;;
LFWG)
  DOWNLOAD_PATH=http://x/cloud-monitor-agent/${TAR_FILE}
  REGION_CODE=cn-beijing
  ;;
HLSY)
  DOWNLOAD_PATH=http://x/cloud-monitor-agent/${TAR_FILE}
  REGION_CODE=cn-beijing
  ;;
BJ)
  DOWNLOAD_PATH=http://x/cloud-monitor-agent/${TAR_FILE}
  REGION_CODE=cn-beijing
  ;;
QD)
  DOWNLOAD_PATH=http://x/cloud-monitor-agent/${TAR_FILE}
  REGION_CODE=cn-qingdao
  ;;
*)
  echo "Unsupported Region"
  exit 1
  ;;
esac

#如果已经安装了监控agent，先卸载旧版agent
WD_SER=$(systemctl status cloud-monitor-agent.service | grep running)

if [[ -n ${WD_SER} ]]; then
  systemctl stop cloud-monitor-agent.service
fi

if [[ -d ${WD_HOME} ]]; then
  rm -rf ${WD_HOME}
  rm -f ${DEST_SERVICE_FILE}
  rm -f /etc/systemd/system/multi-user.target.wants/cloud-monitor-agent.service
fi

download() {
  if [[ -n "${REGION_ID}" ]]; then
    TOS_URL=${DOWNLOAD_PATH}
  else
    echo "unsupported region ${REGION_ID}"
  fi

  echo "download from ${TOS_URL}"

  wget -q "${TOS_URL}" -O "${DEST_TAR_FILE}" -t 3 --connect-timeout=2

  if [[ $? != 0 ]]; then
    echo "download fail, retry..."
    wget -q "${TOS_URL}" -O "${DEST_TAR_FILE}" -t 3 --connect-timeout=2
  fi
}

mkdir -p ${WD_HOME}

if [[ "${WD_OS}" == "linux" ]]; then
  chown -R root:root ${WD_HOME}
fi

download

if [[ ! -f "${DEST_TAR_FILE}" ]]; then
  echo "download failed: {$DEST_TAR_FILE}"
  exit 1
fi

#解压
tar xf "${DEST_TAR_FILE}" -C ${WD_HOME}
rm -f "${DEST_TAR_FILE}"

if [[ -n ${REGION_CODE} ]]; then
  sed -i "s/cn-north-1/${REGION_CODE}/g" ${DEST_MONITOR_CONF}
fi

#拷贝service文件
mv -f ${WD_HOME}/cloud-monitor-agent.service ${DEST_SERVICE_FILE}

#设置watchdog开机自启动
ln -s ${DEST_SERVICE_FILE} /etc/systemd/system/multi-user.target.wants/cloud-monitor-agent.service

#安装状态检测
WD_VERSION=$(${DEST_BIN_FILE} -v)

if [[ -n "${WD_VERSION}" ]]; then
  echo cloud-monitor-agent installed
else
  echo cloud-monitor-agent install failed
  exit 1
fi

systemctl daemon-reload
systemctl start cloud-monitor-agent.service

WD_SER=$(systemctl status cloud-monitor-agent.service | grep running)

if [[ -n ${WD_SER} ]]; then
  echo cloud-monitor-agent start
else
  echo cloud-monitor-agent start failed
  exit 1
fi
