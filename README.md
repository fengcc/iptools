日常开发中经常用到的与IP相关的功能

## 安装
按照[此链接](https://golang.org/doc/install)搭建好Go环境后，在当前目录下执行`go install`即可。

## 功能
### 1. 查询本机外网IP
在终端执行命令`ip`，获得本机外网IP及其相关信息。

![local_ip](img/local_ip.png)

### 2. 查询IP对应信息

如，在终端执行命令`ip 114.114.114.114 `查询对应信息

![ip_find](img/ip_find.png)

### 3. 获取`.ssh/config`中`HostName`字段

例如，`ssh/config`中配置

```
Host test
HostName 123.123.123.123
Port 9999
User root
IdentityFile ~/.ssh/github
```

在终端输入`ip test`即可获得对应IP`123.123.123.123`，同时将此IP复制于剪切板中

![](img/find_ssh.png)

若想要bash能够自动补全`Host`字段，将以下内容添加到`.bashrc`中并执行`source .bashrc`即可。

```shell
[ -e "$HOME/.ssh/config" ] && complete -o "default" -o "nospace" -W "$(grep "^Host " ~/.ssh/config | grep -v "[\d]+\.[\d]+\.[\d]+\.[\d]+" | cut -d " " -f2 | tr ' ' '\n')" ip
```

