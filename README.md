### 模拟信道测速 tc + iperf3
**server client两端都需要安装好`iperf3`, `tc`**

#### 下载软件
[release](https://github.com/lpc-v/netutil/releases)
#### 环境准备
##### server side
```shell
$ iperf3 -s -B yourIP -D                              # yourIP 换成需要测速的ip
$ tc qdisc add dev eth0 root netem delay 10ms loss 1% # eth0 换成需要模拟信道的网口
```

##### client side
```shell
$ tc qdisc add dev eth0 root netem delay 10ms loss 1% # eth0 换成需要模拟信道的网口
```
##### 上传脚本到服务器（server, client, 或者其它能连通都行）
```shell
cd ~ 
mkdir netutil && cd netutil
# ftp上传netutil到/~/netutil
# ftp上传config.ini到/~/netutil
mkdir input && cd input
# ftp上传tc.csv到/~/netutil/input
cd ..
mkdir out
```
##### conf.ini配置文件
![](/img/conf.png)
