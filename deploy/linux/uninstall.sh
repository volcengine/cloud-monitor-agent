#!/bin/bash
#set -x

echo "uninstalling cloud-monitor-agent"

WD_HOME="/usr/local/cloud-monitor-agent"
SERVICE_FILE="/lib/systemd/system/cloud-monitor-agent.service"
LINK_SERVICE_FILE="/etc/systemd/system/multi-user.target.wants/cloud-monitor-agent.service"

CUR_WD=$(systemctl status cloud-monitor-agent.service | grep running)

if [[ -n ${CUR_WD} ]]; then
  systemctl stop cloud-monitor-agent.service
fi

if [[ -d ${WD_HOME} ]]; then
  rm -rf "${WD_HOME}"
  rm -f "${LINK_SERVICE_FILE}"
  rm -f "${SERVICE_FILE}"
fi

echo "cloud-monitor-agent uninstalled"
