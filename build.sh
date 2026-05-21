cd backend
go mod tidy
go build -o ../build/forum .
cd ../frontend
npm install
npm run build
rm -rvf ../build/web
mv dist ../build/web
echo " "
echo "正在运行测试，如果想退出测试，那么请 Ctrl + C"
php -S 0.0.0.0:8000 -t ./build/web