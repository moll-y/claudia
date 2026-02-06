#!/usr/bin/env bash

jq -c '.tool_name //= .notification_type | {tool_name, notification_type}' | timeout 1 socat -u STDIN UNIX-CONNECT:/tmp/claudia.sock
