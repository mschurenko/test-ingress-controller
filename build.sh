#!/usr/bin/env bash

set -e
docker build -t mschurenko/ingress_ctrl .
docker push mschurenko/ingress_ctrl
