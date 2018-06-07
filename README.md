# redid-analyzer
解析rdb, aof, 以及执行monitor, 来查找key和分析各种top-key(big key, hot-key, expiry-key, slowlog-key)


## 安装
go get github.com/lanfang/redis-analyzer

## 如何使用
执行redis-analyzer来查看帮助信息，各功能模块以子命令的方式提供，子命令有自己的参数
```
 redis-analyzer
Usage:
  redis-analyzer [command]

Available Commands:
  bigkey      Find the key over the specified size
  dump        Dump rdb file from redis server
  gen-conf    Generate example json config file
  help        Help about any command
  keys        Grep key through the golang regular
  monitor     A query analyzer that parses redis monitor command
  parse-aof   Parses aof(append-only file) file from local or redis server directly
  parse-rdb   Parses rdb file from local or redis server directly
  slowlog     Collect redis server slowlog

Flags:
  -c, --conf string        json config file
  -e, --etcd-addr string   etcd server address, if not null, we can get config from etcd
  -h, --help               help for redis-analyzer
  -l, --log-file string    log file (default stdout)
  -n, --node strings       redis server master or slave --node "add1, add2, addr"
      --on-master          execute the comand on the master node (default false)
  -o, --output string      the result output file (default stdout)
      --pretty             pretty output result (default false)
  -t, --top-num int        the result top key size (default 10)
      --version            version for redis-analyzer
  -w, --worker-num int     the concurrent worker when get multiple redis server (default 16)

Use "redis-analyzer [command] --help" for more information about a command.
```

### 配置文件
- 支持如下json格式的配置文件, 通过--conf参数指定
```
{
  "nodes": [
    {
      "address": "120.26.10.118:6379", # redis server地址
      "on-master": true # 是否在master 节点执行操作
    }
  ],
  "monitorseconds": 180, # 执行monitor命令的时长(秒)
  "slogloginterval": 10, # redis sloglog监控间隔(秒)
  "rdb-output": "", # 保存rdb文件到指定路径,默认不保存，直连redis解析rdb
  "output": "", # 结果输出文件, 默认标准输出
  "top-num": 10, # 需要获取的top-key的数量
  "bigkey-size": 30,# string类型的big-key大小(KB)
  "element-num": 500000, # 元素个数
  "filter": "", # 查找指定filter的key(golang regexp), 只对keys子命令有效
  "worker-num": 16,# 并发执行的 redis-server
  "log-file": "", # 日志文件, 默认标准输出
  "pretty": true, # 是否格式化输出结果
  "no-expiry": true # 是否只查找未设置超时时间的key, 只对keys子命令有效
}
```

- 通过命令行参数获取配置
命令行参数会覆盖配置json配置文件的配置, 可以通过 --help获取帮助信息:
```
redis-analyzer parse-rdb --help
Parses rdb file from local or redis server directly,
	top-big-key, top-big-type, find-key(abnormaKey key, expiry key and any key you want)

Usage:
  redis-analyzer parse-rdb [flags]

Flags:
  -h, --help                help for parse-rdb
      --rdb-file string     the rdbfile to parse
      --rdb-output string   save rdbfile to rdb-output file

Global Flags:
  -c, --conf string       json config file
  -l, --log-file string   log file (default stdout)
  -n, --node strings      redis server master or slave, support multiple address --node "add1,add2,addr1"
      --on-master         execute the comand on the master node (default false)
  -o, --output string     the result output file (default stdout)
      --pretty            pretty output result (default false)
  -t, --top-num int       the result top key size (default 10)
  -w, --worker-num int    the concurrent worker when get multiple redis server (default 16)
```

### 子命令
#### 解析rdb文件
通过指定rdb文件或者redis服务器，来分析各种key，具体查看帮助信息
```
cat parse_rdb_result_127.0.0.1:6379
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| DATABASE(KEYS)  |          TOP KEY           |  TYPE  | EXIPRY | ELEMENTCOUNT | TOTALSIZE(BYTE) |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| db0:keys=180232 | c89931b27cd527d52:followed | zset   |      0 |        74735 |         3288036 |
+-----------------+----------------------------+--------+        +--------------+-----------------+
| db1:keys=3021   | test_promotion_rule_hash   | hash   |        |          212 |         1959766 |
+-----------------+----------------------------+--------+        +--------------+-----------------+
| db2:keys=104    | test_hot_feed_now_go       | list   |        |          211 |           38352 |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| db3:keys=0      |                            |        |        |              |                 |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| db4:keys=1      | generate_id_test           | string |      0 |            1 |               4 |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| db5:keys=0      |                            |        |        |              |                 |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| db6:keys=14     | logic_worker               | list   |      0 |            1 |             783 |
+-----------------+----------------------------+        +        +--------------+-----------------+
| db7:keys=3      | tj_user_first_time         |        |        |          113 |           28225 |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| db8:keys=0      |                            |        |        |              |                 |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| db9:keys=3      | test_issue_hash            | hash   |      0 |         4212 |          970974 |
+-----------------+----------------------------+        +        +--------------+-----------------+
| db10:keys=499   | lwbie_coupons              |        |        |           40 |            1520 |
+-----------------+----------------------------+--------+        +--------------+-----------------+
| db11:keys=428   | mt_all_total               | string |        |            1 |              75 |
+-----------------+----------------------------+--------+        +--------------+-----------------+
| db12:keys=3292  | morequest_list             | list   |        |         3111 |         4500071 |
+-----------------+----------------------------+--------+        +--------------+-----------------+
| db13:keys=1     | a                          | string |        |            1 |               1 |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| db14:keys=0     |                            |        |        |              |                 |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
| db15:keys=6     | distance:month:2018-05-01  | hash   |      0 |          126 |            6460 |
+-----------------+----------------------------+--------+--------+--------------+-----------------+
database top key
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| DATABASE(KEYS)  |                   TOP KEY                   |  TYPE  |     EXIPRY     | ELEMENTCOUNT | TOTALSIZE(BYTE) |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db0:keys=180232 | to_all_medals                               | string |              0 |            1 |           16865 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | GroupMemberHash                             | hash   |                |          348 |          115850 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | feed_square_feed_go                         | set    |                |         1528 |            9206 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | c89931b27cd527d52:followed                  | zset   |                |        74735 |         3288036 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | geo_common_worker                           | list   |                |         8426 |         1574214 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db1:keys=3021   | sorder:127080397486668622480495             | string |  1530153426032 |            1 |             392 |
+                 +---------------------------------------------+--------+----------------+--------------+-----------------+
|                 | test_promotion_rule_hash                    | hash   |              0 |          212 |         1959766 |
+                 +---------------------------------------------+--------+----------------+--------------+-----------------+
|                 | promotion_fullcut_result                    | set    | 32503651201445 |           29 |             686 |
+                 +---------------------------------------------+--------+----------------+--------------+-----------------+
|                 | cache-engine:tex:exp                        | zset   |              0 |           20 |             743 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | h5game_entity_prx                           | list   |                |            5 |             340 |
+-----------------+---------------------------------------------+--------+                +--------------+-----------------+
| db2:keys=104    | club_total_data_steps:_152                  | string |                |            1 |              12 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | nread_feed_data_go                          | hash   |                |            8 |             130 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | hasclub_microspot_pk_create_tongji_0        | zset   |                |           13 |             572 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | test_hot_feed_now_go                        | list   |                |          211 |           38352 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db3:keys=0      |                                             |        |                |              |                 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db4:keys=1      | generate_id_test                            | string |              0 |            1 |               4 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db5:keys=0      |                                             |        |                |              |                 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db6:keys=14     | _kombu.binding.www_upload_sport_data_worker | set    |              0 |            1 |              60 |
+                 +---------------------------------------------+--------+                +              +-----------------+
|                 | logic_worker                                | list   |                |              |             783 |
+-----------------+---------------------------------------------+--------+                +              +-----------------+
| db7:keys=3      | _kombu.binding.tongji_worker                | set    |                |              |              30 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | tj_user_first_time                          | list   |                |          113 |           28225 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db8:keys=0      |                                             |        |                |              |                 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db9:keys=3      | testissue_index                             | string |              0 |            1 |               4 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | test_issue_hash                             | hash   |                |         4212 |          970974 |
+                 +---------------------------------------------+--------+                +              +-----------------+
|                 | test_issue_set                              | zset   |                |              |           87453 |
+-----------------+---------------------------------------------+--------+                +--------------+-----------------+
| db10:keys=499   | rpt_sync_key                                | string |                |            1 |              17 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | lwbie_coupons                               | hash   |                |           40 |            1520 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | grouprankdaal                               | zset   |                |           10 |             147 |
+-----------------+---------------------------------------------+--------+                +--------------+-----------------+
| db11:keys=428   | mt_all_total                                | string |                |            1 |              75 |
+-----------------+---------------------------------------------+        +                +              +-----------------+
| db12:keys=3292  | test_task_detail_table_all                  |        |                |              |           97133 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | h:generate_cdi                              | hash   |                |       200000 |         2600000 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | fsfg_t                                      | set    |                |         2889 |           16099 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | newsfed_lists                               | zset   |                |          574 |            7528 |
+                 +---------------------------------------------+--------+                +--------------+-----------------+
|                 | morequest_list                              | list   |                |         3111 |         4500071 |
+-----------------+---------------------------------------------+--------+                +--------------+-----------------+
| db13:keys=1     | a                                           | string |                |            1 |               1 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db14:keys=0     |                                             |        |                |              |                 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
| db15:keys=6     | distance:month:2018-05-01                   | hash   |              0 |          126 |            6460 |
+-----------------+---------------------------------------------+--------+----------------+--------------+-----------------+
database top key by type
+-----------------------+----------------------------------------------------------+
|    DATABASE(KEYS)     |                        EXPIRY KEY                        |
+-----------------------+----------------------------------------------------------+
| db0:expiry keys=61569 | 20170418:15-1-7.12.0-19-4.4.2                            |
+-----------------------+----------------------------------------------------------+
| db3:expiry keys=0     |                                                          |
+-----------------------+                                                          +
| db5:expiry keys=0     |                                                          |
+-----------------------+                                                          +
| db8:expiry keys=0     |                                                          |
+-----------------------+----------------------------------------------------------+
| db12:expiry keys=2494 | adv-limit-23-5875958631dfd144233260d58314b689-2017-11-29 |
+-----------------------+----------------------------------------------------------+
| db14:expiry keys=0    |                                                          |
+-----------------------+----------------------------------------------------------+
| db15:expiry keys=6    | distance:week:2018-07-09                                 |
+-----------------------+----------------------------------------------------------+
database expiry key
+----------------------+------------------------------------------------------------+
|    DATABASE(KEYS)    |                        ABNORMAL KEY                        |
+----------------------+------------------------------------------------------------+
| db0:abnormal_keys=6  | bd57e4a5-03a3-4c8e-a9fd-c3375e571f85.2017-07-31            |
|                      | 15:49:44.34176303 +0800 CST                                |
+----------------------+------------------------------------------------------------+
| db1:abnormal_keys=2  | mall_gdpv_90326236074240223121612\u0026pm_r=17053.17030802 |
+----------------------+------------------------------------------------------------+
| db3:abnormal_keys=0  |                                                            |
+----------------------+                                                            +
| db5:abnormal_keys=0  |                                                            |
+----------------------+                                                            +
| db8:abnormal_keys=0  |                                                            |
+----------------------+                                                            +
| db14:abnormal_keys=0 |                                                            |
+----------------------+------------------------------------------------------------+
database abnormal key

database top key: redis服务器里面按照key size排序结果
database top key by type: 每种数据类型按照key排序结果
database expiry key: 根据日期正则匹配，查找未设置过期时间的key(可以通过keys命令实行此功能)
database abnormal key: 根据正则匹配key
```
#### 其它命令
其余命令的使用方法可以参考--help信息
parse-aof只是实现了解析本地aof文件以及连接redis server 作为fake-redis存在, 
未对解析出来的命令做具体的统计操作，如果有需要的话可以自行修改


# Thanks
## - Rdb parse library: [RDB](https://github.com/gohugoio/hugo)
  我把原作者的fork了一份修改为支持list的长度
## - Command line library: [cobra](https://github.com/spf13/cobra)