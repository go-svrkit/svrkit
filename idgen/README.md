# uuid

分布式ID生成


# Usage


### 雪花算法

  时间戳 | WorkerID | SeqID
--------|--------|-------

* 把一个64位的整数分成3个区间，分别用于`时间戳`、`workerID`和`seqID`；
* 时间戳代表时间单元需要一直递增，不同的时间单元实现有毫秒、秒、厘秒等，这里依赖时钟不回调；
* `workerID`用来标识分布式环境下的每个service；
* `seqID`在一个时间单元内持续自增，如果在单个时间单元内`seqID`溢出了，需要sleep等待进入下一个时间单元；

分布式部署下雪花ID只保证了唯一性，无法保证顺序性（递增），
因为ID的分段包含了服务ID，在同个时间单元内(10ms)服务A按时钟后生成的ID会小于服务B按时钟先生成的ID。

不同的雪花算法实现的差异主要集中在3个区间的分配，和workerID自动还是手动分配上。

* sony的实现 `https://github.com/sony/sonyflake`
* 百度的实现 `https://github.com/baidu/uid-generator`
* 美团的实现 `https://github.com/Meituan-Dianping/Leaf`

### 发号器算法

* 算法把一个64位的整数按step划分为N个号段；
* 每个service向发号器申请领取一个代表号段的counter；
* service内部使用这个号段向业务层分配ID；
* service重启或者号段分配完，向发号器申请下一个号段；
* 要求底层存储不能随意被修改影响到上层算法分配；

发号器依赖存储组件，对存储组件的需求是能实现整数自增。

发号器ID只保证了唯一性，无法保证顺序性（递增），因为多个服务同时生成， 如果服务1的生成速度如果比服务2快，服务1的ID号段会先用完，
那么服务1上按时钟先分配的ID会大于服务2上按时钟后分配的ID。

本包提供多种存储选择，redis和mongodb

* redis使用[incr](https://redis.io/commands/incr)命令实现自增
* mongodb使用[findOneAndUpdate](https://docs.mongodb.com/v4.4/reference/method/db.collection.findOneAndUpdate)
