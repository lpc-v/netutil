### netutil 使用说明

>  使用平台：linux, freebsd. 软件在asset/

#### 输入文件

**不要改变文件格式**

` ping.csv` 只需填写IP和时延测试时长即可

`iperf3.csv`A-I列格式固定，可以后面添加备注，对齐即可

#### 使用

- netutil ping

- netutil iperf3

  **server和client上都需要装好iperf3**

#### 输出

/out/ping-*.csv

/out/iperf3-*.csv

