# xxljob_lark_alert
xxljob任务调度只支持邮件告警，懒得重新打包xxljob源码了，随便写一下xxljob的调度失败边缘告警服务，目前只做了飞书告警


docker build -t nexus.kdznwl.cn/aims-test/component/edge-alert-service:1.0.0.0 .

vscode报错：warning: GOPATH set to GOROOT (E:\Code\golang) has no effect
去修改vscode的配置文件 D:\Softs\VSCode-win32-x64-1.68.0\data\user-data\User\settings.json，找到GOROOT和GOPATH的配置，修改为正确的即可

遇到自建包引用报红怎么办，删除引用，保存后重新引用