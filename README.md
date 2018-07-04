# leisure
leisure

# build
./controll.sh build

#dev 环境
# 镜像打包
docker build --build-arg config=./conf/app.dev.ini -t leisure .

# 停止进程
docker stop  leisure && docker rm -f leisure

# docker 启动
docker run --restart=always --name leisure  -idt -p 17005:17005 -p 17006:17006 -v /var/log/ruok/leisure:/home/leisure/logs  leisure:latest 

# 查看日志
tailf /var/log/ruok/leisure/leisure.log