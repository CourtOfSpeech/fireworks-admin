#!/bin/bash

# 租户管理 API 测试脚本
# 使用方式: ./test/api/test_api.sh

BASE_URL="http://localhost:8080/api/v1"

echo "=========================================="
echo "租户管理 API 测试"
echo "=========================================="

echo ""
echo "1. 创建租户 - 企业"
echo "------------------------------------------"
curl -s -X POST "${BASE_URL}/tenants" \
  -H "Content-Type: application/json" \
  -d '{
    "certificate_no": "91110000123456789X",
    "name": "测试公司",
    "type": 1,
    "contact_name": "张三",
    "email": "company@example.com",
    "phone": "13800138001",
    "expired_at": "2026-12-31T23:59:59Z",
    "status": 1
  }'
echo ""

echo ""
echo "2. 创建租户 - 个人"
echo "------------------------------------------"
curl -s -X POST "${BASE_URL}/tenants" \
  -H "Content-Type: application/json" \
  -d '{
    "certificate_no": "110101199001011234",
    "name": "李四",
    "type": 2,
    "contact_name": "李四",
    "email": "personal@example.com",
    "phone": "13800138002",
    "expired_at": "2026-12-31T23:59:59Z",
    "status": 1
  }'
echo ""

echo ""
echo "3. 分页查询租户列表"
echo "------------------------------------------"
curl -s -X GET "${BASE_URL}/tenants?page=1&page_size=10"
echo ""

echo ""
echo "4. 按关键字查询"
echo "------------------------------------------"
curl -s -X GET "${BASE_URL}/tenants?keyword=测试&page=1&page_size=10"
echo ""

echo ""
echo "5. 按状态查询"
echo "------------------------------------------"
curl -s -X GET "${BASE_URL}/tenants?status=1&page=1&page_size=10"
echo ""

echo ""
echo "=========================================="
echo "测试完成"
echo "=========================================="
