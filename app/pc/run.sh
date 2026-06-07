#!/bin/bash
cd "$(dirname "$0")"
echo "校园论坛 PC 客户端"
echo "==================="
echo ""
pip install -r requirements.txt -q
echo ""
python main.py
