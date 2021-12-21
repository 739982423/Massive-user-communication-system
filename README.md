# Massive-user-communication-system
海量用户通讯系统（大概就算个聊天室，用Go实现）
 
### 该仓库为所有源代码，包含：
1. server端代码（server文件夹）
2. client端代码（client文件夹）
3. 一些共用的操作函数，传输格式定义（commom文件夹）

`go.sum`和`go.mod`为golang依赖包管理使用到的文件。使用的Golang版本：1.17。

### 直接使用：
因为服务器端已经部署在阿里云服务器上了，所以可以直接执行目录下的`client.exe`进行聊天室登录。     
服务器的部署使用docker容器进行部署(golang1.17 + redis6.26)，对应的镜像也公开在阿里云。    
拉取指令为：```docker pull registry.cn-beijing.aliyuncs.com/pj_project/go_project:01```


### 自己编译-使用方法：
1. 安装golang的环境：linux下解压官方的.tar.gz文件，将go的解压目录下的/bin文件夹添加入环境变量，再执行source即可
2. 将该项目所有文件放入一个名为finalProject的文件夹内
3. 之后finalProject目录下执行```go build -o client finalProject/client/main```即可对客户端进行编译，将生成client的可执行程序。   
服务器端同理，需要将目录的client改为server。
