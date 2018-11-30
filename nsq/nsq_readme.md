
##
提供了一个最简化的NSQ封装

### 目录结构
```
+-- api
|   +-- nsq_server        //对外总线
|
+-- publisher             //生产方服务
|   +-- publisher         //生产方接口
|   +-- publisher_impl    //生产方接口实现
|
+-- consumer              //消费方服务
|   +-- consumer          //消费方接口
|   +-- consumer_impl     //消费方接口实现
|
+-- message               //对外消息数据结构
|
+-- utils                 //工具包
```