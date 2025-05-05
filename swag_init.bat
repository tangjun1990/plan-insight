echo "格式化swagger注释"
swag fmt
echo "生成 swagger 文档"
swag init -g .\cmd\server\server.go