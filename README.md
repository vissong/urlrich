# urlrich
基于 chromedp&amp;readablity 实现 url rich 化。

## 使用远程 chrome

### 调试

在本地启动，以 mac 为例，先为 chrome 创建一个别名：

```
alias chrome="/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome"
```

然后启动,在办公器需要制定 ioa 代理

```
chrome --headless --remote-debugging-port=9222 --disable-gpu --proxy-server=127.0.0.1:12639
```

启动后会提示在哪个wsurl上启动了浏览器debug，而且会支持http协议用于查看内容，几个已知的：

* http://127.0.0.1:9222/json ：查看已经打开的Tab列表
* http://127.0.0.1:9222/json/version : 查看浏览器版本信息，以及浏览器的debug wsurl，启动的时候提示的那个
* http://127.0.0.1:9222/json/new?http://www.baidu.com : 新开Tab打开指定地址


### 生产