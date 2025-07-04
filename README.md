## 用途
用于批量向不同地址提币，通过**同一网络**提取**同一币种**，但可以对每个地址提取**不同数量**

## 使用方法
1. 去币安开通API，需要绑定IP，开通提币权限

[https://www.binance.com/en/support/faq/detail/360002502072](https://www.binance.com/en/support/faq/detail/360002502072)

[https://developers.binance.com/docs/zh-CN/wallet/general-info](https://developers.binance.com/docs/zh-CN/wallet/general-info)

3. 将向安API的API_KEY和SECRET_KEY填写到.env文件中
3. 在.env文件中配置需要提取的币种以及使用的网络，默认案例为BNB币种和BSC网络
4. 在addresses.csv文件中配置提币目标地址和每个地址的数量，格式在csv文件中
5. 运行可执行文件（可以直接go run main.go）或者下载编译好的执行文件

## 注意事项（主要是网络问题）
1. 国内IP无法裸连币安API，可能需要海外服务器

[https://ipaddress.my/](https://ipaddress.my/)

2. 无法使用梯子的原因：币安API提币需要绑定IP，使用梯子后会改变IP
3. 配置服务器后，币安API需要绑定服务器IP
4. 服务器不能选用美国IP，因为美国是币安的限制地区