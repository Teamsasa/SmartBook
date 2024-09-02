#!/bin/sh
set -e

# マイグレーションを実行
# マイグレーションが成功しなかった場合、再度実行する
until go run internal/migrate/main.go; do
  echo "マイグレーションに失敗しました。再試行します..."
  sleep 1
done
