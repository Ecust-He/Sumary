## SQL优化

### 1  快速学会分析SQL执行效率

#### 定位慢 SQL

当我们实际工作中，碰到某个功能或者某个接口需要很久才能返回结果，我们就应该去确定是不是慢查询导致的。定位慢 SQL 有如下两种解决方案：

- 查看慢查询日志确定已经执行完的慢查询
- show processlist 查看正在执行的慢查询

#####  通过慢查询日志

默认情况下，慢查询日志中不会记录管理语句，可通过设置 log_slow_admin_statements = on 让管理语句中的慢查询也会记录到慢查询日志中。

默认情况下，也不会记录查询时间不超过 long_query_time 但是不使用索引的语句，可通过配置 log_queries_not_using_indexes = on 让不使用索引的 SQL 都被记录到慢查询日志中（即使查询时间没超过 long_query_time 配置的值）。

如果需要使用慢查询日志，一般分为四步：**开启慢查询日志**、**设置慢查询阀值**、**确定慢查询日志路径**、**确定慢查询日志的文件名**。

###### 开启慢查询

```mysql
mysql> set global slow_query_log = on;

Query OK, 0 rows affected (0.00 sec)
```

###### 设置慢查询时间阀值

```mysql
mysql> set global long_query_time = 1;

Query OK, 0 rows affected (0.00 sec)
```

###### 确定慢查询日志路径

```mysql
mysql> show global variables like "datadir";

+---------------+------------------------+
| Variable_name | Value                  |
+---------------+------------------------+
| datadir       | /data/mysql/data/3306/ |
+---------------+------------------------+

1 row in set (0.00 sec)
```

###### 确定慢查询日志的文件名

```mysql
mysql> show global variables like "slow_query_log_file";

+---------------------+----------------+
| Variable_name       | Value          |
+---------------------+----------------+
| slow_query_log_file | mysql-slow.log |
+---------------------+----------------+

1 row in set (0.00 sec)
```

根据上面的查询结果，可以直接查看 /data/mysql/data/3306/mysql-slow.log 文件获取已经执行完的慢查询

```mysql
[root@mysqltest ~]# tail -n5 /data/mysql/data/3306/mysql-slow.log

Time: 2019-05-21T09:15:06.255554+08:00

User@Host: root[root] @ localhost []  Id: 8591152

Query_time: 10.000260  Lock_time: 0.000000 Rows_sent: 1  Rows_examined: 0

SET timestamp=1558401306;
select sleep(10);
```

这里对上方的执行结果详细描述一下：

- tail -n5：只查看慢查询文件的最后 5 行
- Time：慢查询发生的时间
- User@Host：客户端用户和 IP
- Query_time：查询时间
- Lock_time：等待表锁的时间
- Rows_sent：语句返回的行数
- Rows_examined：语句执行期间从存储引擎读取的行数

上面这种方式是用系统自带的慢查询日志查看的，如果觉得系统自带的慢查询日志不方便查看，小伙伴们可以使用 pt-query-digest 或者 mysqldumpslow 等工具对慢查询日志进行分析。

##### 通过 show processlist

有时慢查询正在执行，已经导致数据库负载偏高了，而由于慢查询还没执行完，因此慢查询日志还看不到任何语句。此时可以使用 show processlist 命令判断正在执行的慢查询。show processlist 显示哪些线程正在运行。如果有 PROCESS 权限，则可以看到所有线程。否则，只能看到当前会话的线程。

```mysql
mysql> show processlist\G`

`......`

`*************************** 10. row ***************************`

     `Id: 7651833`

   `User: one`

   `Host: 192.168.1.251:52154`

     `db: ops`

`Command: Query`

   `Time: 3`

  `State: User sleep`

   `Info: select sleep(10)`

`......`

`10 rows in set (0.00 sec)`
```

- Time：表示执行时间
- Info：表示 SQL 语句

可以通过它的执行时间（Time）来判断是否是慢 SQL。

#### 使用 explain 分析慢查询

分析 SQL 执行效率是优化 SQL 的重要手段，通过上面讲的两种方法，定位到慢查询语句后，我们就要开始分析 SQL 执行效率了，子曾经曰过：“工欲善其事，必先利其器”，我们可以通过 explain、show profile 和 trace 等诊断工具来分析慢查询。

**Explain 可以获取 MySQL 中 SQL 语句的执行计划，比如语句是否使用了关联查询、是否使用了索引、扫描行数等。可以帮我们选择更好地索引和写出更优的 SQL 。**

使用方法：在查询语句前面加上 explain 运行就可以了。

为了便于理解，先创建两张测试表（方便第 1、2 节实验使用），建表及数据写入语句如下：

```mysql
CREATE DATABASE muke;           /* 创建测试使用的database，名为muke */
use muke;                       /* 使用muke这个database */
drop table if exists t1;        /* 如果表t1存在则删除表t1 */

CREATE TABLE `t1` (             /* 创建表t1 */
  `id` int(11) NOT NULL auto_increment,
  `a` int(11) DEFAULT NULL,
  `b` int(11) DEFAULT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_a` (`a`),
  KEY `idx_b` (`b`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;	

drop procedure if exists insert_t1; /* 如果存在存储过程insert_t1，则删除 */
delimiter ;;
create procedure insert_t1()        /* 创建存储过程insert_t1 */
begin
  declare i int;                    /* 声明变量i */
  set i=1;                          /* 设置i的初始值为1 */
  while(i<=1000)do                  /* 对满足i<=1000的值进行while循环 */
    insert into t1(a,b) values(i, i); /* 写入表t1中a、b两个字段，值都为i当前的值 */
    set i=i+1;                      /* 将i加1 */
  end while;
end;;
delimiter ;                 /* 创建批量写入1000条数据到表t1的存储过程insert_t1 */
call insert_t1();           /* 运行存储过程insert_t1 */

drop table if exists t2;    /* 如果表t2存在则删除表t2 */
create table t2 like t1;    /* 创建表t2，表结构与t1一致 */
insert into t2 select * from t1;   /* 将表t1的数据导入到t2 */
```

试使用 explain 分析一条 SQL，例子如下：

```mysql
mysql> explain select * from t1 where b=100;
```

##### Explain 的结果各字段

| 列名            | 解释                                                         |
| :-------------- | :----------------------------------------------------------- |
| id              | 查询编号                                                     |
| **select_type** | 查询类型：显示本行是简单还是复杂查询                         |
| table           | 涉及到的表                                                   |
| partitions      | 匹配的分区：查询将匹配记录所在的分区。仅当使用 partition 关键字时才显示该列。对于非分区表，该值为 NULL。 |
| **type**        | 本次查询的表连接类型                                         |
| possible_keys   | 可能选择的索引                                               |
| **key**         | 实际选择的索引                                               |
| key_len         | 被选择的索引长度：一般用于判断联合索引有多少列被选择了       |
| ref             | 与索引比较的列                                               |
| **rows**        | 预计需要扫描的行数，对 InnoDB 来说，这个值是估值，并不一定准确 |
| filtered        | 按条件筛选的行的百分比                                       |
| **Extra**       | 附加信息                                                     |

######  select_type

| select_type 的值     | 解释                                                         |
| :------------------- | :----------------------------------------------------------- |
| SIMPLE               | 简单查询 (不使用关联查询或子查询)                            |
| PRIMARY              | 如果包含关联查询或者子查询，则最外层的查询部分标记为 primary |
| UNION                | 联合查询中第二个及后面的查询                                 |
| DEPENDENT UNION      | 满足依赖外部的关联查询中第二个及以后的查询                   |
| UNION RESULT         | 联合查询的结果                                               |
| SUBQUERY             | 子查询中的第一个查询                                         |
| DEPENDENT SUBQUERY   | 子查询中的第一个查询，并且依赖外部查询                       |
| DERIVED              | 用到派生表的查询                                             |
| MATERIALIZED         | 被物化的子查询                                               |
| UNCACHEABLE SUBQUERY | 一个子查询的结果不能被缓存，必须重新评估外层查询的每一行     |
| UNCACHEABLE UNION    | 关联查询第二个或后面的语句属于不可缓存的子查询               |

###### type

| type 的值       | 解释                                                         |
| :-------------- | :----------------------------------------------------------- |
| system          | 查询对象表只有一行数据，且只能用于 MyISAM 和 Memory 引擎的表，这是最好的情况 |
| const           | 基于主键或唯一索引查询，最多返回一条结果                     |
| eq_ref          | 表连接时基于主键或非 NULL 的唯一索引完成扫描                 |
| ref             | 基于普通索引的等值查询，或者表间等值连接                     |
| fulltext        | 全文检索                                                     |
| ref_or_null     | 表连接类型是 ref，但进行扫描的索引列中可能包含 NULL 值       |
| index_merge     | 利用多个索引                                                 |
| unique_subquery | 子查询中使用唯一索引                                         |
| index_subquery  | 子查询中使用普通索引                                         |
| range           | 利用索引进行范围查询                                         |
| index           | 全索引扫描                                                   |
| ALL             | 全表扫描                                                     |

**上表的这些情况，查询性能从上到下依次是最好到最差。**

###### Extra

| Extra 常见的值                        | 解释                                                         | 例子                                                         |
| :------------------------------------ | :----------------------------------------------------------- | :----------------------------------------------------------- |
| Using filesort                        | 将用外部排序而不是索引排序，数据较小时从内存排序，否则需要在磁盘完成排序 | explain select * from t1 order by create_time;               |
| Using temporary                       | 需要创建一个临时表来存储结构，通常发生对没有索引的列进行 GROUP BY 时 | explain select * from t1 group by create_time;               |
| Using index                           | 使用覆盖索引                                                 | explain select a from t1 where a=111;                        |
| Using where                           | 使用 where 语句来处理结果                                    | explain select * from t1 where create_time=‘2019-06-18 14:38:24’; |
| Impossible WHERE                      | 对 where 子句判断的结果总是 false 而不能选择任何数据         | explain select * from t1 where 1<0;                          |
| Using join buffer (Block Nested Loop) | 关联查询中，被驱动表的关联字段没索引                         | explain select * from t1 straight_join t2 on (t1.create_time=t2.create_time); |
| Using index condition                 | 先条件过滤索引，再查数据                                     | explain select * from t1 where a >900 and a like “%9”;       |
| Select tables optimized away          | 使用某些聚合函数（比如 max、min）来访问存在索引的某个字段是  | explain select max(a) from t1;                               |

#### 使用show profile 分析慢查询

**有时需要确定 SQL 到底慢在哪个环节，此时 explain 可能不好确定。**在 MySQL 数据库中，通过 profile，能够更清楚地了解 SQL 执行过程的资源使用情况，能让我们知道到底慢在哪个环节。

可以通过配置参数 profiling = 1 来启用 SQL 分析。该参数可以在全局和 session 级别来设置。对于全局级别则作用于整个MySQL 实例，而 session 级别仅影响当前 session 。该参数开启后，后续**执行的 SQL 语句都将记录其资源开销，如 IO、上下文切换、CPU、Memory**等等。根据这些开销进一步分析当前 SQL 从而进行优化与调整。

##### 确定是否支持 profile

```mysql
mysql> select @@have_profiling;

+------------------+
| @@have_profiling |
+------------------+
| YES              |
+------------------+

1 row in set, 1 warning (0.00 sec)
```

##### 查看 profiling 是否关闭的

```mysql
mysql> select @@profiling;

+-------------+
| @@profiling |
+-------------+
|           0 |
+-------------+

1 row in set, 1 warning (0.00 sec)
```

##### 通过 set 开启 profile

```mysql
mysql> set profiling=1;

Query OK, 0 rows affected, 1 warning (0.00 sec)
```

##### 执行SQL语句

```mysql
mysql> select * from t1 where b=1000;
```

##### 确定 SQL 的 query id

```mysql
mysql> show profiles;
+----------+------------+-------------------------------+
| Query_ID | Duration   | Query                         |
+----------+------------+-------------------------------+
|        1 | 0.00063825 | select * from t1 where b=1000 |
+----------+------------+-------------------------------+
1 row in set, 1 warning (0.00 sec)
```

##### 查询 SQL 执行详情

通过 show profile for query 可看到执行过的 SQL 每个状态和消耗时间：

```mysql
mysql> show profile for query 1;
+----------------------+----------+
| Status               | Duration |
+----------------------+----------+
| starting             | 0.000115 |
| checking permissions | 0.000013 |
| Opening tables       | 0.000027 |
| init                 | 0.000035 |
| System lock          | 0.000017 |
| optimizing           | 0.000016 |
| statistics           | 0.000025 |
| preparing            | 0.000020 |
| executing            | 0.000006 |
| Sending data         | 0.000294 |
| end                  | 0.000009 |
| query end            | 0.000012 |
| closing tables       | 0.000011 |
| freeing items        | 0.000024 |
| cleaning up          | 0.000016 |
+----------------------+----------+
15 rows in set, 1 warning (0.00 sec)
```

#### trace 分析 SQL 优化器

从前面学到了 explain 可以查看 SQL 执行计划，但是无法知道它为什么做这个决策，如果想确定多种索引方案之间是如何选择的或者排序时选择的是哪种排序模式，有什么好的办法吗？

通过trace，能够进一步了解为什么优化器选择A执行计划而不是选择B执行计划，或者知道某个排序使用的排序模式，帮助我们更好地理解优化器行为。

```mysql
select * from t1 where a >900 and b > 910 order  by a;
```

通过上面执行计划中 key 这个字段可以看出，该语句使用的是 b 字段的索引 idx_b。实际表 t1 中，a、b 两个字段都有索引，为什么条件中有这两个索引字段却偏偏选了 b 字段的索引呢？这时就可以使用 trace 进行分析。大致步骤如下：

```mysql
mysql> set session optimizer_trace="enabled=on",end_markers_in_json=on;
/* optimizer_trace="enabled=on" 表示开启 trace；end_markers_in_json=on 表示 JSON 输出开启结束标记 */
Query OK, 0 rows affected (0.00 sec)

mysql> select * from t1 where a >900 and b > 910 order  by a;
+------+------+------+
| id   | a    | b    |
+------+------+------+
|    1 |    1 |    1 |
|    2 |    2 |    2 |

......

| 1000 | 1000 | 1000 |
+------+------+------+
1000 rows in set (0.00 sec)

mysql> SELECT * FROM information_schema.OPTIMIZER_TRACE\G
*************************** 1. row ***************************
QUERY: select * from t1 where a >900 and b > 910 order  by a    --SQL语句
TRACE: {
  "steps": [
    {
      "join_preparation": {				--SQL准备阶段
        "select#": 1,
        "steps": [
          {
            "expanded_query": "/* select#1 */ select `t1`.`id` AS `id`,`t1`.`a` AS `a`,`t1`.`b` AS `b`,`t1`.`create_time` AS `create_time`,`t1`.`update_time` AS `update_time` from `t1` where ((`t1`.`a` > 900) and (`t1`.`b` > 910)) order by `t1`.`a`"
          }
        ] /* steps */
      } /* join_preparation */
    },
    {
      "join_optimization": {			--SQL优化阶段
        "select#": 1,
        "steps": [
          {
            "condition_processing": {    --条件处理
              "condition": "WHERE",
              "original_condition": "((`t1`.`a` > 900) and (`t1`.`b` > 910))",        --原始条件
              "steps": [
                {
                  "transformation": "equality_propagation",
                  "resulting_condition": "((`t1`.`a` > 900) and (`t1`.`b` > 910))" 		--等值传递转换
                },
                {
                  "transformation": "constant_propagation",
                  "resulting_condition": "((`t1`.`a` > 900) and (`t1`.`b` > 910))"       --常量传递转换
                },
                {
                  "transformation": "trivial_condition_removal",
                  "resulting_condition": "((`t1`.`a` > 900) and (`t1`.`b` > 910))"        --去除没有的条件后的结构
                }
              ] /* steps */
            } /* condition_processing */
          },
          {
            "substitute_generated_columns": {
            } /* substitute_generated_columns */   --替换虚拟生成列
          },
          {
            "table_dependencies": [		--表依赖详情
              {
                "table": "`t1`",
                "row_may_be_null": false,
                "map_bit": 0,
                "depends_on_map_bits": [
                ] /* depends_on_map_bits */
              }
            ] /* table_dependencies */
          },
          {
            "ref_optimizer_key_uses": [
            ] /* ref_optimizer_key_uses */
          },
          {
            "rows_estimation": [	--预估表的访问成本
              {
                "table": "`t1`",
                "range_analysis": {
                  "table_scan": {
                    "rows": 1000,       --扫描行数
                    "cost": 207.1       --成本
                  } /* table_scan */,
                  "potential_range_indexes": [    --分析可能使用的索引
                    {
                      "index": "PRIMARY",
                      "usable": false,       --为false，说明主键索引不可用
                      "cause": "not_applicable"
                    },
                    {
                      "index": "idx_a",      --可能使用索引idx_a
                      "usable": true,
                      "key_parts": [
                        "a",
                        "id"
                      ] /* key_parts */
                    },
                    {
                      "index": "idx_b",      --可能使用索引idx_b
                      "usable": true,
                      "key_parts": [
                        "b",
                        "id"
                      ] /* key_parts */
                    }
                  ] /* potential_range_indexes */,
                  "setup_range_conditions": [
                  ] /* setup_range_conditions */,
                  "group_index_range": {
                    "chosen": false,
                    "cause": "not_group_by_or_distinct"
                  } /* group_index_range */,
                  "analyzing_range_alternatives": { --分析各索引的成本
                    "range_scan_alternatives": [
                      {
                        "index": "idx_a",	--使用索引idx_a的成本
                        "ranges": [
                          "900 < a"			--使用索引idx_a的范围
                        ] /* ranges */,
                        "index_dives_for_eq_ranges": true, --是否使用index dive（详细描述请看下方的知识扩展）
                        "rowid_ordered": false, --使用该索引获取的记录是否按照主键排序
                        "using_mrr": false,  	--是否使用mrr
                        "index_only": false,    --是否使用覆盖索引
                        "rows": 100,            --使用该索引获取的记录数
                        "cost": 121.01,         --使用该索引的成本
                        "chosen": true          --可能选择该索引
                      },
                      {
                        "index": "idx_b",       --使用索引idx_b的成本
                        "ranges": [
                          "910 < b"
                        ] /* ranges */,
                        "index_dives_for_eq_ranges": true,
                        "rowid_ordered": false,
                        "using_mrr": false,
                        "index_only": false,
                        "rows": 90,
                        "cost": 109.01,
                        "chosen": true             --也可能选择该索引
                      }
                    ] /* range_scan_alternatives */,
                    "analyzing_roworder_intersect": { --分析使用索引合并的成本
                      "usable": false,
                      "cause": "too_few_roworder_scans"
                    } /* analyzing_roworder_intersect */
                  } /* analyzing_range_alternatives */,
                  "chosen_range_access_summary": {  --确认最优方法
                    "range_access_plan": {
                      "type": "range_scan",
                      "index": "idx_b",
                      "rows": 90,
                      "ranges": [
                        "910 < b"
                      ] /* ranges */
                    } /* range_access_plan */,
                    "rows_for_plan": 90,
                    "cost_for_plan": 109.01,
                    "chosen": true
                  } /* chosen_range_access_summary */
                } /* range_analysis */
              }
            ] /* rows_estimation */
          },
          {
            "considered_execution_plans": [  --考虑的执行计划
              {
                "plan_prefix": [
                ] /* plan_prefix */,
                "table": "`t1`",
                "best_access_path": {          --最优的访问路径
                  "considered_access_paths": [ --决定的访问路径
                    {
                      "rows_to_scan": 90,      --扫描的行数
                      "access_type": "range",  --访问类型：为range
                      "range_details": {
                        "used_index": "idx_b"  --使用的索引为：idx_b
                      } /* range_details */,
                      "resulting_rows": 90,    --结果行数
                      "cost": 127.01,          --成本
                      "chosen": true,		   --确定选择
                      "use_tmp_table": true
                    }
                  ] /* considered_access_paths */
                } /* best_access_path */,
                "condition_filtering_pct": 100,
                "rows_for_plan": 90,
                "cost_for_plan": 127.01,
                "sort_cost": 90,
                "new_cost_for_plan": 217.01,
                "chosen": true
              }
            ] /* considered_execution_plans */
          },
          {
            "attaching_conditions_to_tables": {  --尝试添加一些其他的查询条件
              "original_condition": "((`t1`.`a` > 900) and (`t1`.`b` > 910))",
              "attached_conditions_computation": [
              ] /* attached_conditions_computation */,
              "attached_conditions_summary": [
                {
                  "table": "`t1`",
                  "attached": "((`t1`.`a` > 900) and (`t1`.`b` > 910))"
                }
              ] /* attached_conditions_summary */
            } /* attaching_conditions_to_tables */
          },
          {
            "clause_processing": {
              "clause": "ORDER BY",
              "original_clause": "`t1`.`a`",
              "items": [
                {
                  "item": "`t1`.`a`"
                }
              ] /* items */,
              "resulting_clause_is_simple": true,
              "resulting_clause": "`t1`.`a`"
            } /* clause_processing */
          },
          {
            "reconsidering_access_paths_for_index_ordering": {
              "clause": "ORDER BY",
              "index_order_summary": {
                "table": "`t1`",
                "index_provides_order": false,
                "order_direction": "undefined",
                "index": "idx_b",
                "plan_changed": false
              } /* index_order_summary */
            } /* reconsidering_access_paths_for_index_ordering */
          },
          {
            "refine_plan": [          --改进的执行计划
              {
                "table": "`t1`",
                "pushed_index_condition": "(`t1`.`b` > 910)",
                "table_condition_attached": "(`t1`.`a` > 900)"
              }
            ] /* refine_plan */
          }
        ] /* steps */
      } /* join_optimization */
    },
    {
      "join_execution": {             --SQL执行阶段
        "select#": 1,
        "steps": [
          {
            "filesort_information": [
              {
                "direction": "asc",
                "table": "`t1`",
                "field": "a"
              }
            ] /* filesort_information */,
            "filesort_priority_queue_optimization": {
              "usable": false,             --未使用优先队列优化排序
              "cause": "not applicable (no LIMIT)"     --未使用优先队列排序的原因是没有limit
            } /* filesort_priority_queue_optimization */,
            "filesort_execution": [
            ] /* filesort_execution */,
            "filesort_summary": {           --排序详情
              "rows": 90,
              "examined_rows": 90,          --参与排序的行数
              "number_of_tmp_files": 0,     --排序过程中使用的临时文件数
              "sort_buffer_size": 115056,
              "sort_mode": "<sort_key, additional_fields>"   --排序模式（详解请看下方知识扩展）
            } /* filesort_summary */
          }
        ] /* steps */
      } /* join_execution */
    }
  ] /* steps */
}
MISSING_BYTES_BEYOND_MAX_MEM_SIZE: 0	--该字段表示分析过程丢弃的文本字节大小，本例为0，说明没丢弃任何文本
          INSUFFICIENT_PRIVILEGES: 0    --查看trace的权限是否不足，0表示有权限查看trace详情
1 row in set (0.00 sec)
------------------------------------------------
------------------------------------------------


mysql> set session optimizer_trace="enabled=off";
/* 及时关闭trace */
```

TRACE 字段中整个文本大致分为三个过程。

- 准备阶段：对应文本中的 join_preparation
- 优化阶段：对应文本中的 join_optimization
- 执行阶段：对应文本中的 join_execution

使用时，重点关注优化阶段和执行阶段。

由此例可以看出：

- 在 trace 结果的 analyzing_range_alternatives 这一项可以看到：使用索引 idx_a 的成本为 121.01，使用索引 idx_b 的成本为 109.01，显然使用索引 idx_b 的成本要低些，因此优化器选择了 idx_b 索引；
- 在 trace 结果的 filesort_summary 这一项可以看到：排序模式为<sort_key, additional_fields>，表示使用的是单路排序，即一次性取出满足条件行的所有字段，然后在 sort buffer 中进行排序。

### 2   条件字段有索引，为什么查询也这么慢?

如果我们想在某一本书中找到特定的主题，一般最快的方法是先看索引，找到对应的主题在哪个页码。

而对于 MySQL 而言，如果需要查找某一行的值，可以先通过索引找到对应的值，然后根据索引匹配的记录找到需要查询的数据行。

#### 函数操作

##### 验证对条件字段做函数操作是否能走索引

首先创建测试表，建表及数据写入语句如下：

```mysql
use muke;                       /* 使用muke这个database */

drop table if exists t1;        /* 如果表t1存在则删除表t1 */

CREATE TABLE `t1` (             /* 创建表t1 */
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `a` varchar(20) DEFAULT NULL,
  `b` int(20) DEFAULT NULL,
  `c` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,  
  PRIMARY KEY (`id`),
  KEY `idx_a` (`a`) USING BTREE,
  KEY `idx_b` (`b`) USING BTREE,
  KEY `idx_c` (`c`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

drop procedure if exists insert_t1; /* 如果存在存储过程insert_t1，则删除 */
delimiter ;;
create procedure insert_t1()        /* 创建存储过程insert_t1 */
begin
  declare i int;                    /* 声明变量i */
  set i=1;                          /* 设置i的初始值为1 */
  while(i<=10000)do                 /* 对满足i<=10000的值进行while循环 */
    insert into t1(a,b) values(i,i);  /* 写入表t1中a、b两个字段，值都为i当前的值 */
    set i=i+1;                        /* 将i加1 */
  end while;
end;;
delimiter ;
call insert_t1();                    /* 运行存储过程insert_t1 */

update t1 set c = '2019-05-22 00:00:00';  /* 更新表t1的c字段，值都为'2019-05-22 00:00:00' */
update t1 set c = '2019-05-21 00:00:00' where id=10000;	 /* 将id为10000的行的c字段改为与其它行都不一样的数据，以便后面实验使用 */
```

对于上面创建的测试表，比如要查询测试表 t1 单独某一天的所有数据，SQL如下：

```mysql
select * from t1 where date(c) ='2019-05-21';
```

查看图中的执行计划，type 为 ALL，key 字段结果为 NULL，因此知道该 SQL 是没走索引的全表扫描。

##### 对条件字段做函数操作不走索引的原因

索引树中存储的是列的实际值和主键值。如果拿 ‘2019-05-21’ 去匹配，将无法定位到索引树中的值。因此放弃走索引，而选择全表扫描。

##### 函数操作的 SQL 优化

因此如果需要优化的话，改成 c 字段实际值相匹配的形式。因为 SQL 的目的是查询 2019-05-21 当天所有的记录，因此可以改成范围查询，如下：

```mysql
select * from t1 where c>='2019-05-21 00:00:00' and c<='2019-05-21 23:59:59';
```

根据上面的结果，可确定，走了 c 字段的索引（对应关注字段 key），扫描行数 1 行（对应关注字段 rows）。

#### 隐式转换

##### 验证隐式转换是否能走索引

比如我们要查询 a 字段等于 1000 的值，SQL如下：

```mysql
mysql> select * from t1 where a=1000;
+------+------+------+---------------------+
| id   | a    | b    | c                   |
+------+------+------+---------------------+
| 1000 | 1000 | 1000 | 2019-05-22 00:00:00 |
+------+------+------+---------------------+
1 row in set (0.00 sec)
```

通过 type 这列可以看到是最差的情况 ALL（全表扫描，如果对 explain 结果中 type 各个值没印象的，可以查看第 2 节中<表 3-type 各项值解释>）， 通过 key 这列可以看到没走 a 字段的索引，通过 rows 这列可以看到进行了全表扫描。

##### 不走索引的原因

a 字段类型是 varchar(20)，而语句中 a 字段条件值没加单引号，导致 MySQL 内部会先把a转换成int型，再去做判断，相当于实际执行的 SQL 语句如下：

```mysql
mysql> select * from t1 where cast(a as signed int) =1000;
```

##### 隐式转换的 SQL 优化

```mysql
mysql> explain select * from t1 where a='1000';
```

通过 type 这列，可以看到是 ref（基于普通索引的等值查询，比 ALL 性能好很多，可复习第 2 节<表3-type 各项值解释>），通过key这列，可以看到已经走了 a 字段的索引，通过rows这列可以看到通过索引查询后就扫描了一行。

####  模糊查询

#####  分析模糊查询

很多时候我们想根据某个字段的某几个关键字查询数据，比如会有如下 SQL：

```mysql
mysql> select * from t1 where a like '%1111%';
+------+------+------+---------------------+
| id   | a    | b    | c                   |
+------+------+------+---------------------+
| 1111 | 1111 | 1111 | 2019-05-22 00:00:00 |
+------+------+------+---------------------+
1 row in set (0.00 sec)
```

重点留意type、key、rows、Extra，发现是全表扫描。

##### 模糊查询优化建议

改业务，让模糊查询必须包含条件字段前面的值，然后落到数据库的查询为：

```mysql
mysql> select * from t1 where a like '1111%';
+------+------+------+---------------------+
| id   | a    | b    | c                   |
+------+------+------+---------------------+
| 1111 | 1111 | 1111 | 2019-05-22 00:00:00 |
+------+------+------+---------------------+
1 row in set (0.00 sec)
```

#### 范围查询

##### 构造不能使用索引的范围查询

我们拿测试表举例，比如要取出b字段1到2000范围数据，SQL 如下 ：

```mysql
mysql> select * from t1 where b>=1 and b <=2000;
```

发现并不能走b字段的索引。

原因：优化器会根据检索比例、表大小、I/O块大小等进行评估是否使用索引。比如单次查询的数据量过大，优化器将不走索引。

##### 优化范围查询

降低单次查询范围，分多次查询：

```mysql
mysql> select * from t1 where b>=1 and b <=1000;
mysql> select * from t1 where b>=1001 and b <=2000;
```

实际这种范围查询而导致使用不了索引的场景经常出现，比如按照时间段抽取全量数据，每条SQL抽取一个月的；或者某张业务表历史数据的删除。遇到此类操作时，**应该在执行之前对SQL做explain分析，确定能走索引，再进行操作**，否则不但可能导致操作缓慢，在做更新或者删除时，甚至会导致表所有记录锁住，十分危险。

#### 计算操作

##### 查询条件进行计算操作的 SQL 执行效率

```mysql
mysql> explain select * from t1 where b-1 =1000;
```

原因：对索引字段做运算将使用不了索引。

经验分享：**一般需要对条件字段做计算时，建议通过程序代码实现，而不是通过MySQL实现。如果在MySQL中计算的情况避免不了，那必须把计算放在等号后面。**

### 3  如何优化数据导入？

我们有时会遇到批量数据导入的场景，而数据量稍微大点，会发现导入非常耗时间。

#### 一次插入多行值

插入行所需的时间由以下因素决定

- 连接：30%
- 向服务器发送查询：20%
- 解析查询：20%
- 插入行：10% * 行的大小
- 插入索引：10% * 索引数
- 结束：10%

##### 导出一条 SQL 包含多行数据的数据文件

```mysql
[root@mysqltest muke]# mysqldump -utest_user3 -p'userBcdQ19Ic' -h127.0.0.1 --set-gtid-purged=off --single-transaction --skip-add-locks  muke t1 >t1.sql
```

这里对上面 mysqldump 所使用到的一些参数做下解释：

| 参数                | 详解                                                         |
| :------------------ | :----------------------------------------------------------- |
| -utest_user3        | 用户名，这里使用的是root用户                                 |
| -p’userBcdQ19Ic’    | 密码                                                         |
| -h127.0.0.1         | 连接的MySQL服务端IP                                          |
| set-gtid-purged=off | 不添加SET @@GLOBAL.GTID_PURGED                               |
| –single-transaction | 设置事务的隔离级别为可重复读，即REPEATABLE READ，这样能保证在一个事务中所有相同的查询读取到同样的数据 |
| –skip-add-locks     | 取消每个表导出之前加lock tables操作。                        |
| muke                | 库名                                                         |
| t1                  | 表名                                                         |
| t1.sql              | 导出数据到这个文件                                           |

查看文件 t1.sql 内容，可看到数据是单条 SQL 有多条数据，如下：

```mysql
......
DROP TABLE IF EXISTS `t1`;		
/* 按照上面的备份语句备份的数据文件包含drop命令时，需要特别小心，在后续使用备份文件做导入操作时，应该确定所有表名，防止drop掉业务正在使用的表 */
......
CREATE TABLE `t1`......
......
INSERT INTO `t1` VALUES (1,'1',1,'2019-05-24 15:44:10'),(2,'2',2,'2019-05-24 15:44:10'),(3,'3',3,'2019-05-24 15:44:10')......
......
```

##### 导出一条SQL只包含一行数据的数据文件

```mysql
[root@mysqltest muke]# mysqldump -utest_user3 -p'userBcdQ19Ic' -h127.0.0.1 --set-gtid-purged=off --single-transaction --skip-add-locks --skip-extended-insert muke t1 >t1_row.sql
```

mysqldump命令参数解释：

| 参数                  | 详解                          |
| :-------------------- | :---------------------------- |
| –skip-extended-insert | 一条SQL一行数据的形式导出数据 |

备份文件t1_row.sql内容如下

```mysql
......
INSERT INTO `t1` VALUES (1,'1',1,'2019-05-24 15:44:10');
INSERT INTO `t1` VALUES (2,'2',2,'2019-05-24 15:44:10');
INSERT INTO `t1` VALUES (3,'3',3,'2019-05-24 15:44:10');
......
```

#####  导入时间的对比

首先导入一行 SQL 包含多行数据的数据文件：

```sql
[root@mysqltest ~]# time mysql -utest_user3 -p'userBcdQ19Ic' -h127.0.0.1 muke <t1.sql

real	0m0.230s
user	0m0.007s
sys		0m0.003s
```

耗时0.2秒左右。

导入一条SQL只包含一行数据的数据文件：

```sql
[root@mysqltest ~]# time mysql -utest_user3 -p'userBcdQ19Ic' -h127.0.0.1 muke <t1_row.sql

real	0m31.138s
user	0m0.088s
sys		0m0.126s
```

耗时31秒左右。

##### 结论

```mysql
一次插入多行花费时间0.2秒，一次插入一行花费了31秒，对比效果明显，因此建议有大批量导入时，推荐一条insert语句插入多行数据。
```

#### 关闭自动提交

##### 对比开启和关闭自动提交的效率

Autocommit 开启时会为每个插入执行提交。可以在InnoDB导入数据时，关闭自动提交。如下：

```mysql
SET autocommit=0;
INSERT INTO `t1` VALUES (1,'1',1,'2019-05-24 15:44:10');
INSERT INTO `t1` VALUES (2,'2',2,'2019-05-24 15:44:10');
INSERT INTO `t1` VALUES (3,'3',3,'2019-05-24 15:44:10');
......
COMMIT;
```

开启自动提交的情况下导入是31秒。

**关闭自动提交的情况下导入是1秒左右，因此导入多条数据时，关闭自动提交，让多条 insert 一次提交，可以大大提升导入速度。**

##### 原因分析

与本节前面讲的一次插入多行能提高批量插入速度的原因一样，因为批量导入大部分时间耗费在客户端与服务端通信的时间，所以多条 insert 语句合并提交可以减少客户端与服务端通信的时间，并且合并提交还可以减少数据落盘的次数。

#### 参数调整

影响MySQL写入速度的主要两个参数：innodb_flush_log_at_trx_commit、sync_binlog。

### 4  让order by、group by查询更快

#### order by 原理

#####  MySQL 的排序方式

按照排序原理分，MySQL 排序方式分两种：

- 通过有序索引直接返回有序数据
- 通过 Filesort 进行的排序

##### 怎么确定某条排序的 SQL 所使用的排序方式？

使用 explain 来查看该排序 SQL 的执行计划，重点关注 Extra 字段：

如果该字段里显示是 Using index，则表示是通过有序索引直接返回有序数据。

如果该字段里显示是 Using filesort，则表示该 SQL 是通过 Filesort 进行的排序。

##### Filesort 是在内存中还是在磁盘中完成排序的？

MySQL 中的 Filesort 并不一定是在磁盘文件中进行排序的，也有可能在内存中排序，内存排序还是磁盘排序取决于排序的数据大小和 sort_buffer_size 配置的大小。

- 如果 “排序的数据大小” < sort_buffer_size: 内存排序
- 如果 “排序的数据大小” > sort_buffer_size: 磁盘排序

#####  Filesort 下的排序模式

Filesort 下的排序模式有三种，具体介绍如下：

- < sort_key, rowid >双路排序（又叫回表排序模式）：是首先根据相应的条件取出相应的排序字段和可以直接定位行数据的行 ID，然后在 sort buffer 中进行排序，排序完后需要再次取回其它需要的字段；
- < sort_key, additional_fields >单路排序：是一次性取出满足条件行的所有字段，然后在sort buffer中进行排序；
- < sort_key, packed_additional_fields >打包数据排序模式：与单路排序相似，区别是将 char 和 varchar 字段存到 sort buffer 中时，更加紧缩。

#### order by 优化

##### 添加合适索引

######  排序字段添加索引

首先我们看下对 d 字段（没有索引）进行排序的执行计划：

```sql
explain select d,id from t1 order by d;
```

发现使用的是 filesort（关注 Extra 字段）。

再看些对 c 字段（有索引）进行排序的执行计划：

```sql
explain select c,id from t1 order by c;
```

可以看到，根据有索引的字段排序，在 Extra 中显示的就为 Using index，表示使用的是索引排序。如果数据量比较大，显然通过有序索引直接返回有序数据效率更高。

**因此可以在排序字段上添加索引来优化排序语句。**

###### 多个字段排序优化

有时面对的需求是要对多个字段进行排序，而这种情况应该怎么优化或者设计索引呢？

对 a、c 两个字段进行排序的执行计划：

```sql
explain select id,a,c from t1 order by a,c;
```

观察 Extra 字段，发现使用的是 filesort。

再看对 a、b（a、b 两个字段有联合索引）两个字段进行排序：

```sql
explain select id,a,b from t1 order by a,b;
```

发现使用的是索引排序。

多个字段排序的情况，如果要通过添加索引优化，得注意排序字段的顺序与联合索引中列的顺序要一致。

**如果多个字段排序，可以在多个排序字段上添加联合索引来优化排序语句。**

###### 先等值查询再排序的优化

我们更多的情况是会先根据某个字段条件查出一部分数据，然后再排序，而这类 SQL 应该如果优化呢？

表 t1中，根据 a=1000 过滤数据再根据 d 字段排序的执行计划如下：

```sql
explain select id,a,d from t1 where a=1000 order by d;
```

可以在 Extra 字段中看到 “Using filesort”，说明使用的是 filesort 排序。

再看下根据 a=1000 过滤数据在根据 b 字段排序的执行计划（a、b 两个字段有联合索引）：

```sql
explain select id,a,b from t1 where a=1000 order by b;
```

可以在 Extra 字段中看到“Using index”，说明使用的是索引排序。

因此，对于先等值查询再排序的语句，可以通过在**条件字段**和**排序字段**添加**联合索引**来优化此类排序语句。

##### 去掉不必要的返回字段

根据执行计划的结果，可以看到，查询所有字段的这条 SQL 是 filesort 排序，而只查 id、a、b 三个字段的 SQL 是 index 排序，为什么查询所有字段会不走索引？

这个例子中，查询所有字段不走索引的原因是：扫描整个索引并查找到没索引的行的成本比扫描全表的成本更高，所以优化器放弃使用索引。

##### 修改参数

##### 几种无法利用索引排序的情况

###### 使用范围查询再排序

对于先等值过滤再排序的语句，可以通过在条件字段和排序字段添加联合索引来优化；但是如果联合索引中前面的字段使用了范围查询，对后面的字段排序是否能用到索引排序呢？下面我们通过实验验证一下：

```sql
explain select id,a,b from t1 where a>9000 order by b;
```

这里对上面执行计划做下解释：首先条件 a>9000 使用了索引（关注 key 字段对应的值为 idx_a_b）；在 Extra 中，看到了“Using filesort”，表示使用了 filesort 排序，并没有使用索引排序。所以联合索引中前面的字段使用了范围查询，对后面的字段排序使用不了索引排序。

原因是：a、b 两个字段的联合索引，对于单个 a 的值，b 是有序的。而对于 a 字段的范围查询，也就是 a 字段会有多个值，取到 a，b 的值 b 就不一定有序了，因此要额外进行排序。

###### ASC 和 DESC 混合使用将无法使用索引

对联合索引多个字段同时排序时，如果一个是顺序，一个是倒序，则使用不了索引，如下例：

```sql
explain select id,a,b from t1 order by a asc,b desc;
```

####  group by 优化

默认情况，会对 group by 字段排序，因此优化方式与 order by 基本一致，如果目的只是分组而不用排序，可以指定 order by null 禁止排序。

### 5  换种思路写分页查询

#### 根据自增且连续主键排序的分页查询

首先来看一个根据自增且连续主键排序的分页查询的例子：

```sql
select * from t1 limit 99000,2;
```

该 SQL 表示查询从第 99001开始的两行数据，没添加单独 order by，表示通过主键排序。我们再看表 t1，因为主键是自增并且连续的，所以可以改写成按照主键去查询从第 99001开始的两行数据，如下：

```sql
select * from t1 where id >99000 limit 2;
```

查询的结果是一致的。我们再对比一下执行计划：
原 SQL 中 key 字段为 NULL，表示未走索引，rows 显示 99965，表示扫描的行数 99965行；

改写后的 SQL key 字段为 PRIMARY，表示走了主键索引，扫描了1000行。

显然改写后的 SQL 执行效率更高。

另外如果原 SQL 是 order by 非主键的字段，按照上面说的方法改写会导致两条 SQL 的结果不一致。所以这种改写得满足以下两个条件：

- 主键自增且连续
- 结果是按照主键排序的

#### 查询根据非主键字段排序的分页查询

再看一个根据非主键字段排序的分页查询，SQL 如下：

```sql
select * from t1 order by a limit 99000,2;
```

查询时间是 0.08 秒。

我们来看下这条 SQL 的执行计划：
发现并没有使用 a 字段的索引（key 字段对应的值为 null），具体原因：**扫描整个索引并查找到没索引的行的成本比扫描全表的成本更高，所以优化器放弃使用索引**。

知道不走索引的原因，那么怎么优化呢？

其实关键是**让排序时返回的字段尽可能少**，所以可以让排序和分页操作先查出主键，然后根据主键查到对应的记录，SQL 改写如下（这里参考了《深入浅出 MySQL》18.4.7 优化分页查询）：

```sql
select * from t1 f inner join (select id from t1 order by a limit 99000,2)g on f.id = g.id;
```

需要的结果与原 SQL 一致，执行时间 0.02 秒，是原 SQL 执行时间的四分之一，我们再对比优化前后的执行计划：
原 SQL 使用的是 filesort 排序，而优化后的 SQL 使用的是索引排序。

### 6  Join语句可以这样优化

#### 关联查询的算法

#####  Nested-Loop Join 算法

一个简单的 Nested-Loop Join(NLJ) 算法一次一行循环地从第一张表（称为驱动表）中读取行，在这行数据中取到关联字段，根据关联字段在另一张表（被驱动表）里取出满足条件的行，然后取出两张表的结果合集。

我们试想一下，如果在被驱动表中这个关联字段没有索引，那么每次取出驱动表的关联字段在被驱动表查找对应的数据时，都会对被驱动表做一次全表扫描，成本是非常高的（比如驱动表数据量是 m，被驱动表数据量是 n，则扫描行数为 m * n ）。

好在 MySQL 在关联字段有索引时，才会使用 NLJ，如果没索引，就会使用 Block Nested-Loop Join，等下会细说这个算法。我们先来看下在有索引情况的情况下，使用 Nested-Loop Join 的场景（称为：Index Nested-Loop Join）。

因为 MySQL 在关联字段有索引时，才会使用 NLJ，因此本节后面的内容所用到的 NLJ 都表示 Index Nested-Loop Join。

如下例：

```sql
select * from t1 inner join t2 on t1.a = t2.a;       /* sql1 */
```

> Tips：表 t1 和表 t2 中的 a 字段都有索引。

怎么确定这条 SQL 使用的是 NLJ 算法？

我们先来看下 sql1 的执行计划：
从执行计划中可以看到这些信息：

- 驱动表是 t2，被驱动表是 t1。原因是：explain 分析 join 语句时，在第一行的就是驱动表；选择 t2 做驱动表的原因：如果没固定连接方式（比如没加 straight_join）优化器会优先选择小表做驱动表。**所以使用 inner join 时，前面的表并不一定就是驱动表。**
- 使用了 NLJ。原因是：一般 join 语句中，如果执行计划 Extra 中未出现 Using join buffer （***）；则表示使用的 join 算法是 NLJ。

sql1 的大致流程如下：

1. 从表 t2 中读取一行数据；
2. 从第 1 步的数据中，取出关联字段 a，到表 t1 中查找；
3. 取出表 t1 中满足条件的行，跟 t2 中获取到的结果合并，作为结果返回给客户端；
4. 重复上面 3 步。

在这个过程中会读取 t2 表的所有数据，因此这里扫描了 100 行，然后遍历这 100 行数据中字段 a 的值，根据 t2 表中 a 的值索引扫描 t1 表中的对应行，这里也扫描了 100 行。因此整个过程扫描了 200 行。

##### Block Nested-Loop Join 算法

我们一起看看下面这条 SQL 语句：

```sql
select * from t1 inner join t2 on t1.b = t2.b;       /* sql2 */
```

> Tips：表 t1 和表 t2 中的 b 字段都没有索引
> 在 Extra 发现 Using join buffer (Block Nested Loop)，这个就说明该关联查询使用的是 BNL 算法。

我们再看下 sql2 的执行流程：

1. 把 t2 的所有数据放入到 join_buffer 中
2. 把表 t1 中每一行取出来，跟 join_buffer 中的数据做对比
3. 返回满足 join 条件的数据

在这个过程中，对表 t1 和 t2 都做了一次全表扫描，因此扫描的总行数为10000(表 t1 的数据总量) + 100(表 t2 的数据总量) = 10100。并且 join_buffer 里的数据是无序的，因此对表 t1 中的每一行，都要做 100 次判断，所以内存中的判断次数是 100 * 10000= 100 万次。

下面我们来回答上面提出的一个问题：

如果被驱动表的关联字段没索引，为什么会选择使用 BNL 算法而不继续使用 Nested-Loop Join 呢？

在被驱动表的关联字段没索引的情况下，比如 sql2：

如果使用 Nested-Loop Join，那么扫描行数为 100 * 10000 = 100万次，这个是磁盘扫描。

如果使用 BNL，那么磁盘扫描是 100 + 10000=10100 次，在内存中判断 100 * 10000 = 100万次。

显然后者磁盘扫描的次数少很多，因此是更优的选择。因此对于 MySQL 的关联查询，如果被驱动表的关联字段没索引，会使用 BNL 算法。

##### Batched Key Access 算法

#### 小表做驱动表

**建议小表驱动大表**

####  临时表

多数情况我们可以通过在被驱动表的关联字段上加索引来让 join 使用 NLJ 或者 BKA，但有时因为某条关联查询只是临时查一次，如果再去添加索引可能会浪费资源，那么有什么办法优化呢？

**当遇到 BNL 的 join 语句，如果不方便在关联字段上添加索引，不妨尝试创建临时表，然后在临时表中的关联字段上添加索引，然后通过临时表来做关联查询。**
