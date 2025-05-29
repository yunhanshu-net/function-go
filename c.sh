git init
git remote add origin https://github.com/yunhanshu-net/function-go.git
git add .
git commit -m "feat: 初始化 function-go SDK

- 重命名 sdk-go 为 function-go
- 保持原有功能和API不变
- 支持云函数开发和运行时集成"
git branch -M main
git push -u origin main
cd ..