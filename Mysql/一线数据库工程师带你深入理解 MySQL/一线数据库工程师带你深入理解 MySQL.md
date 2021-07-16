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

### 7  为何count（*）这么慢？

#### 重新认识 count()

#####  count(a) 和 count(*) 的区别

当 count() 统计某一列时，比如 count(a)，a 表示列名，是不统计 null 的。

而 count(*) 无论是否包含空值，都会统计。

#####  MyISAM 引擎和 InnoDB 引擎 count(*) 的区别

对于 MyISAM 引擎，如果没有 where 子句，也没检索其它列，那么 count(*) 将会非常快。因为 MyISAM 引擎会把表的总行数存在磁盘上。

首先我们看下对 t2 表（存储引擎为 MyISAM）不带 where 子句做 count(*) 的执行计划：

```sql
explain select count(*) from t2;
```

在 Extra 字段发现 “Select tables optimized away” 关键字，表示是从 MyISAM 引擎维护的准确行数上获取到的统计值。

而 InnoDB 并不会保留表中的行数，**因为并发事务可能同时读取到不同的行数**。所以执行 count(*) 时都是临时去计算的，会比 MyISAM 引擎慢很多。

对比 MyISAM 引擎和 InnoDB 引擎 count(*) 的区别：

- MyISAM 会维护表的总行数，放在磁盘中，如果有 count(*) 的需求，直接返回这个数据
- 但是 InnoDB 就会去遍历普通索引树，计算表数据总量

##### MySQL 5.7.18 前后 count(*) 的区别

在 MySQL 5.7.18 之前，InnoDB 通过扫描聚簇索引来处理 count(*) 语句。

从 MySQL 5.7.18 开始，通过遍历最小的可用二级索引来处理 count(*) 语句。如果不存在二级索引，则扫描聚簇索引。但是，如果索引记录不完全在缓存池中的话，处理 count(*) 也是比较久的。

新版本为什么会使用二级索引来处理 count(*) 语句呢？

原因是 InnoDB 二级索引树的叶子节点上存放的是主键，而主键索引树的叶子节点上存放的是整行数据，所以二级索引树比主键索引树小。因此优化器基于成本的考虑，优先选择的是二级索引。所以 count(主键) 其实没 count (*) 快。

##### count(1) 比 count(*) 快吗？

count(*) 无论是否包含空值，所有结果都会统计。

而 count(1)中的 1 是恒真表达式，因此也会统计所有结果。

所以 count(1) 和 count(*) 统计结果没差别。

#### 哪些方法可以加快 count()？

##### show table status

```sql
show table status like 't1';
```

##### 用 Redis 做计数器

##### 增加计数表

## MYSQL索引

### 1  为什么添加索引能提高查询速度？

#### 跟索引相关的一些算法

对于 MySQL 而言，使用最频繁的就是 B+ 树索引，所以我们必须要知道 B+ 树的结构，而 B+ 树是借鉴了二分查找法、二叉查找树、平衡二叉树、B 树的一些思想构建的。

##### 二分查找法

二分查找法的查找过程是：将记录按顺序排列，查找时先以有序列的中点位置为比较对象，如果要找的元素值小于该中点元素，则将查询范围缩小为左半部分；如果要找的元素值大于该中点元素，则将查询范围缩小为右半部分。以此类推，直到查到需要的值。

##### 二叉查找树

二分查找法的查找过程是：将记录按顺序排列，查找时先以有序列的中点位置为比较对象，如果要找的元素值小于该中点元素，则将查询范围缩小为左半部分；如果要找的元素值大于该中点元素，则将查询范围缩小为右半部分。以此类推，直到查到需要的值。

##### 平衡二叉树

我们一起看下平衡二叉树的定义：满足二叉查找树的定义，另外必须满足任何节点的两个子树的高度差最大为 1。

##### B 树

B 树可以理解为一个节点可以拥有多于 2 个子节点的多叉查找树。

B 树中同一键值不会出现多次，要么在叶子节点，要么在内节点上。

##### B+ 树

B+ 树是 B 树的变体，定义基本与 B 树一致，与 B 树的不同点：

- 所有叶子节点中包含了全部关键字的信息
- 各叶子节点用指针进行连接
- 非叶子节点上只存储 key 的信息，这样相对 B 树，可以增加每一页中存储 key 的数量。
- B 树是纵向扩展，最终变成一个 “瘦高个”，而 B+ 树是横向扩展的

B+树与B树的结构最大的区别就是：

它的键一定会出现在叶子节点上，同时也有可能在非叶子节点中重复出现。而 B 树中同一键值不会出现多次。

#### B+ 树索引

B+ 树索引就是基于本节前面介绍的 B+ 树发展而来的。在数据库中，B+ 树的高度一般都在 2 ~ 4 层，所以**查找某一行数据最多只需要 2 到 4 次 IO。而没索引的情况，需要逐行扫描，明显效率低很多，这也就是为什么添加索引能提高查询速度。**

B+ 树索引并不能找到一个给定键值的具体行，B+ 树索引能找到的只是被查找数据行所在的页。然后数据库通过把页读入到缓冲池（buffer pool）中，在内存中通过二分查找法进行查找，得到需要的数据。

InnoDB 中 B+ 树索引分为聚集索引和辅助索引，我们再继续了解这两种索引的特点。

#####  聚集索引

InnoDB 的数据是按照主键顺序存放的，而聚集索引就是按照每张表的主键构造一颗 B+ 树，它的叶子节点存放的是整行数据。

InnoDB 的主键一定是聚集索引。如果没有定义主键，聚集索引可能是第一个不允许为 null 的唯一索引，也有可能是 row id。

由于实际的数据页只能按照一颗 B+ 树进行排序，因此每张表只能有一个聚集索引（TokuDB 引擎除外）。查询优化器倾向于采用聚集索引，因为聚集索引能够在 B+ 树索引的叶子节点上直接找到数据。

聚集索引对于主键的排序查找和范围查找速度非常快。

两点关键信息：

- 根据主键值创建了 B+ 树结构
- 每个叶子节点包含了整行数据

##### 辅助索引

上图中两点关键点需要注意：

- 根据 a 字段的值创建了 B+ 树结构
- 每个叶子节点保存的是 a 字段自己的键值和主键 ID

下面这条查询语句：

```sql
select * from t8 where a=3;
```

它先通过 a 字段上的索引树，得到主键 id 为 3，再到 id 的聚集索引树上找到对应的行数据。

而下面这条 SQL:

```sql
select * from t8 where id=3;
```

查询到的结果是一样的，而执行过程则只需要搜索 id 的聚集索引树。我们能看出辅助索引的查询比主键查询多扫描一颗索引树，所以，我们应该**尽量使用主键做为条件进行查询**。

### 2  哪些情况下需要添加索引？

目前比较常见需要创建索引的场景有：数据检索时在条件字段添加索引、聚合函数对聚合字段添加索引、对排序字段添加索引、为了防止回表添加索引、关联查询在关联字段添加索引等。

#### 数据检索

把有索引的字段 a 作为条件进行查询

```sql
select * from t9_1 where a = 90000;
```

执行计划中：

前者 type 字段为 ALL，后者 type 字段为 ref，显然后者性能更好

rows 这个字段前者是 1596288，而后者是 16，有索引的情况扫描行数大大降低。

因此建议数据检索时，在条件字段添加索引。

#### 聚合函数

求有索引的字段 a 的最大值：

```sql
select max(a) from t9_1;
```

相比对没有索引的字段 d 求最大值（花费330毫秒），**显然索引能提升 max() 函数的效率，同理也能提升 min() 函数的效率**。

从 MySQL 5.7.18 开始，通过遍历最小的可用二级索引来处理 count(*) 语句，如果不存在二级索引，则扫描聚簇索引。原因是：InnoDB 二级索引树的叶子节点上存放的是主键，而主键索引树的叶子节点上存放的是整行数据，所以二级索引树比主键索引树小。因此优化器基于成本的考虑，优先选择的是二级索引。

**因此索引对聚合函数 count(\*) 也有优化作用。**

#### 排序

通过添加合适索引优化 order by 的方法

- 如果对单个字段排序，则可以在这个排序字段上添加索引来优化排序语句；
- 如果是多个字段排序，可以在多个排序字段上添加联合索引来优化排序语句；
- 如果是先等值查询再排序，可以通过在条件字段和排序字段添加联合索引来优化排序语句。

####  避免回表

比如下面这条 SQL：

```sql
select a,d from t9_1 where a=90000;
```

可以走 a 字段的索引，但是在学了第 8 节后，我们知道了辅助索引的结构，如果通过辅助索引来寻找数据，InnoDB 存储引擎会遍历辅助索引树查找到对应记录的主键，然后通过主键索引回表去找对应的行数据。

但是，如果条件字段和需要查询的字段有联合索引的话，其实回表这一步就省了，因为联合索引中包含了这两个字段的值。像这种索引就已经覆盖了我们的查询需求的场景，我们称为：**覆盖索引**。比如下面这条 SQL：

```sql
select b,c from t9_1 where b=90000;
预览
```

可直接通过联合索引 idx_b_c 找到 b、c 的值

**所以可以通过添加覆盖索引让 SQL 不需要回表，从而减少树的搜索次数，让查询更快地返回结果。**

####  关联查询

### 3   普通索引和唯一索引有哪些区别？

对于普通索引和唯一索引的区别，也许你已经知道：有普通索引的字段可以写入重复的值，而有唯一索引的字段不可以写入重复的值。其实对于 MySQL 来说，不止这一种区别。

#### Insert Buffer

为什么要增加 Insert Buffer？

增加 Insert Buffer 有两个好处：

- 减少磁盘的离散读取
- 将多次插入合并为一次操作

但是得注意的是，使用 Insert Buffer 得满足两个条件：

- 索引是辅助索引
- 索引不是唯一

跟 Insert Buffer 一样，Change Buffer 也得满足这两个条件：

- 索引是辅助索引
- 索引不是唯一

#### Change Buffer

为什么唯一索引的更新不使用 Change Buffer ?

原因：唯一索引**必须要将数据页读入内存才能判断是否违反唯一性约束**。如果都已经读入到内存了，那直接更新内存会更快，就没必要使用 Change Buffer 了。

#### 普通索引和唯一索引的区别

- 数据修改时，普通索引可以用 Change Buffer，而唯一索引不行。
- 数据修改时，唯一索引在 RR 隔离级别下，更容易出现死锁。
- 查询数据是，普通索引查到满足条件的第一条记录还需要继续查找下一个记录，而唯一索引查找到第一个记录就可以直接返回结果了，但是普通索引多出的查找次数所消耗的资源多数情况可以忽略不计。

####  普通索引和唯一索引如何选择

如果业务要求某个字段唯一，但是代码不能完全保证写入唯一值，则添加唯一索引，让这个字段唯一，该字段新增重复数据时，将报类似如下的错：

```sql
ERROR 1062 (23000): Duplicate entry '1' for key 'f1'
```

如果代码确定某个字段不会有重复的数据写入，则可以选择添加普通索引。 因为普通索引可以使用 Change Buffer，并且出现死锁的概率比唯一索引低。

### 4  联合索引有哪些讲究？

#### 认识联合索引

联合索引：是指对表上的多个列进行索引。适合 where 条件中的多列组合，在某些场景可以避免回表。

联合索引的键值数量大于 1（比如上图中有 a 和 b 两个键值），与单个键值的 B+ 树一样，也是按照键值排序的。比如图中 a、b 两个字段的值为 (1,1),(1,2),(1,3),(2,1),(2,2),(2,3)，是按(a,b) 进行排序的。因此，**对于 a、b 两个字段都做为条件时，查询是可以走索引的；对于单独 a 字段查询也是可以走索引的。但是对于 b 字段单独查询就走不了索引了**，因为在上图，b 字段对应的值为 1,2,3,1,2,3，显然不是有序的，所以走不了 b 字段的索引。

联合索引的建议：

- where 条件中，经常同时出现的列放在联合索引中。
- 把选择性最大的列放在联合索引的最左边。

#### 联合索引使用分析

##### 可以完整用到联合索引的情况

下面我们列出几种可以完整用到联合索引的情况

```sql
select * from t11 where a=1 and b=1 and c=1; /* sql1 */
```

##### 只能使用部分联合索引的情况

有些场景只能用到部分联合索引，这里就列出几种情况。

```sql
select * from t11 where a=1 and b=1; /* sql11 */
```

对于联合索引 idx_a_b_c（a,b,c） ，如果条件中只包含 a 和 c，则只能用到联合索引中 a 的索引。c 这里是用不了索引的。**联合索引 idx_a_b_c(a,b,c) 相当于 (a) 、(a,b) 、(a,b,c) 三种索引，称为联合索引的最左原则。**

##### 可以用到覆盖索引的情况

什么是覆盖索引？

从辅助索引中就可以查询到结果，不需要回表查询聚集索引中的记录。

使用覆盖索引的优势：因为不需要扫描聚集索引，因此可以减少 SQL 执行过程的 IO 次数。

```sql
select b,c from t11 where a=3; /* sql21 */
```

通过 a 字段上的条件，去联合索引 idx_a_b_c 的索引树上可以直接查找到 b 字段和 c 字段的值，不需要回表，因此 sql21 使用到了覆盖索引。

通过 a 字段上的条件，去联合索引 idx_a_b_c 的索引树上可以直接查找到 b 字段和 c 字段的值，不需要回表，因此 sql21 使用到了覆盖索引。

```sql
select c from t11 where a=1 and b=1 ; /* sql22 */
```

跟 sql21 类似，在联合索引 idx_a_b_c 的索引树上，通过 a 和 b 的值可以直接找到 c 的值，因此 sql22 使用的也是覆盖索引。

```sql
select id from t11 where a=1 and b=1 and c=1; /* sql23 */
```

通过 a、b、c 三个字段的值，去联合索引树的叶子节点找到主键 id，不需要回表，因此 sql23 也使用了覆盖索引。

### 5  为什么MySQL会选错索引？

#### show index的使用

上面几个重要的字段做一下解释：

- Non_unique：如果是唯一索引，则值为 0，如果可以有重复值，则值为 1
- Key_name：索引名字
- Seq_in_index：索引中的列序号，比如联合索引 idx_a_b_c (a,b,c) ，那么三个字段分别对应 1,2,3
- Column_name：字段名
- Collation：字段在索引中的排序方式，A 表示升序，NULL 表示未排序
- Cardinality：索引中不重复记录数量的预估值，该值等会儿会详细讲解
- Sub_part：如果是前缀索引，则会显示索引字符的数量；如果是对整列进行索引，则该字段值为 NULL
- Null：如果列可能包含空值，则该字段为 YES；如果不包含空值，则该字段值为 ’ ’
- Index_type：索引类型，包括 BTREE、FULLTEXT、HASH、RTREE 等

#### Cardinality 取值

#### 统计信息不准确导致选错索引

在 MySQL 中，**优化器控制着索引的选择。一般情况下，优化器会考虑扫描行数、是否使用临时表、是否排序等因素，然后选择一个最优方案去执行 SQL 语句**。

#### 单次选取的数据量过大导致选错索引

## MySQL锁

### 1  全局锁和表锁什么场景会用到？

**如何保证数据并发访问的一致性、有效性呢？**

MySQL 中，锁就是协调多个用户或者客户端并发访问某一资源的机制，保证数据并发访问时的一致性和有效性。

根据加锁的范围，MySQL 中的锁可分为三类：

- 全局锁
- 表级锁
- 行锁

#### 全局锁

MySQL 全局锁会关闭所有打开的表，并使用全局读锁锁定所有表。其命令为：

```sql
FLUSH TABLES WITH READ LOCK;
```

简称：FTWRL，可以使用下面命令解锁：

```sql
UNLOCK TABLES;
```

进行 FTWRL 实验：

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| FLUSH TABLES WITH READ LOCK; Query OK, 0 rows affected (0.00 sec) |                                                              |
| select * from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** | select * from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** |
| insert into t14(a,b) values(2,2); ERROR 1223 (HY000): Can’t execute the query because you have a conflicting read lock **（报错）** | insert into t14(a,b) values(2,2);/* sql1 */ **（等待）**     |
| UNLOCK TABLES;                                               | insert into t14(a,b) values(2,2);/* sql1 */ Query OK, 1 row affected (5.73 sec) **（session1 解锁后，在等待的 sql1 马上执行成功）** |

上面的实验中，当 session1 执行 FTWRL 后，本线程 session1 和其它线程 session2 都可以查询，本线程和其它线程都不能更新。

原因是：**当执行 FTWRL 后，所有的表都变成只读状态，数据更新或者字段更新将会被阻塞。**

全局锁一般用在整个库（包含非事务引擎表）做备份（mysqldump 或者 xtrabackup）时。也就是说，在整个备份过程中，整个库都是只读的，其实这样风险挺大的。如果是在主库备份，会导致业务不能修改数据；而如果是在从库备份，就会导致主从延迟。

好在 mysqldump 包含一个参数 --single-transaction，可以在一个事务中创建一致性快照，然后进行所有表的备份。因此增加这个参数的情况下，备份期间可以进行数据修改。但是需要所有表都是事务引擎表。所以这也是建议使用 InnoDB 存储引擎的原因之一。

而对于 xtrabackup，可以分开备份 InnoDB 和 MyISAM，或者不执行 --master-data，可以避免使用全局锁。

#### 表级锁

表级锁有两种：表锁和元数据锁。

##### 表锁

表锁使用场景：

1. 事务需要更新某张大表的大部分或全部数据。如果使用默认的行锁，不仅事务执行效率低，而且可能造成其它事务长时间锁等待和锁冲突，这种情况下可以考虑使用表锁来提高事务执行速度；
2. 事务涉及多个表，比较复杂，可能会引起死锁，导致大量事务回滚，可以考虑表锁避免死锁。

其中表锁又分为表读锁和表写锁，命令分别是：

表读锁：

```sql
lock tables t14 read;
```

表写锁：

```sql
lock tables t14  write;
```

下面我们分别用实验验证表读锁和表写锁。

表读锁实验：

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| lock tables t14 read; Query OK, 0 rows affected (0.00 sec)   |                                                              |
| select id,a,b from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** | select id,a,b from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** |
| insert into t14(a,b) values(3,3); ERROR 1099 (HY000): Table ‘t14’ was locked with a READ lock and can’t be updated **（报错）** | insert into t14(a,b) values(3,3);/* sql2 */ **（等待）**     |
| unlock tables; Query OK, 0 rows affected (0.00 sec)          | insert into t14(a,b) values(3,3);/* sql2 */ Query OK, 1 row affected (10.97 sec) **（session1 解锁后，sql2 立马写入成功）** |

从上面的实验我们可以看出，在 session1 中对表 t14 加表读锁，session1 和 session2 都可以查询表 t14 的数据；而 session1 执行更新会报错，session2 执行更新会等待（直到 session1 解锁后才更新成功）。

总结：**对表执行 lock tables xxx read （表读锁）时，本线程和其它线程可以读，本线程写会报错，其它线程写会等待。**

我们再来看一下表写锁实验：

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| lock tables t14 write; Query OK, 0 rows affected (0.00 sec)  |                                                              |
| select id,a,b from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** | select id,a,b from t14 limit 1;/* sql3 */ **（等待）**       |
| unlock tables; Query OK, 0 rows affected (0.01 sec)          | select id,a,b from t14 limit 1;/* sql3 */ … 1 row in set (7.16 sec) **（session1 解锁后，sql3 马上返回查询结果）** |
| lock tables t14 write; Query OK, 0 rows affected (0.00 sec)  |                                                              |
| delete from t14 limit 1; Query OK, 1 row affected, 1 warning (0.00 sec) **（能正常执行删除语句）** | delete from t14 limit 1;/* sql4 */ **（等待）**              |
| unlock tables; Query OK, 0 rows affected (0.00 sec)          | delete from t14 limit 1;/* sql4 */ Query OK, 1 row affected, 1 warning (14.94 sec) **（session1 解锁后，sql4 立马执行成功）** |

总结：**对表执行 lock tables xxx write （表写锁）时，本线程可以读写，其它线程读写都会阻塞。**

##### 元数据锁

#### 总结

其中**全局锁会让所有的表变成只读状态，所有更新操作都会被阻塞。**

而表级锁分为表锁和元数据锁。

表锁又提到了表读锁和表写锁，并都进行了实验。两者的区别是：

**表读锁：本线程和其它线程可以读，本线程写会报错，其它线程写会等待。**

**表写锁：本线程可以读写，其它线程读写都会阻塞。**

为了保证事务和 DDl 并行执行数据一致，在 MySQL 5.5.3 引入了 MDL 锁。通过本节讲解的 MDL 锁机制，应该注意的几个点是：

- 尽量避免慢查询
- 事务要及时提交
- 避免大事务
- 避免在业务高峰执行 DDL 操作

### 2  行锁：InnoDB替代MyISAM的重要原因

MySQL 5.5 之前的默认存储引擎是 MyISAM，5.5 之后改成了 InnoDB。InnoDB 后来居上最主要的原因就是：

- InnoDB 支持事务：适合在并发条件下要求数据一致的场景。
- InnoDB 支持行锁：有效降低由于删除或者更新导致的锁定。

#### 两阶段锁

传统的关系型数据库加锁的一个原则是：两阶段锁原则。

两阶段锁：锁操作分为两个阶段，加锁阶段和解锁阶段，并且保证加锁阶段和解锁阶段不相交。

我们可以通过下面这张表理解两阶段锁：

| 序号 | MySQL 操作      | 解释                                    | 锁阶段   |
| :--- | :-------------- | :-------------------------------------- | :------- |
| 1    | begin;          | 事务开始                                |          |
| 2    | insert into …;  | 加 insert 对应的锁                      | 加锁阶段 |
| 3    | update table …; | 加 update 对应的锁                      | 加锁阶段 |
| 4    | delete from …;  | 加 delete 对应的锁                      | 加锁阶段 |
| 5    | commit;         | 事务结束，同时释放 2、3、4 步骤中加的锁 | 解锁阶段 |

####  InnoDB 行锁模式

InnoDB 实现了以下两种类型的行锁：

- 共享锁（S）：允许一个事务去读一行，阻止其它事务获得相同数据集的排他锁；
- 排他锁（X）：允许获得排他锁的事务更新数据，阻止其它事务取得相同数据集的共享读锁和排他写锁。

对于普通 select 语句，InnoDB 不会加任何锁，事务可以通过以下语句显式给记录集加共享锁或排他锁：

- 共享锁（S）：select * from table_name where … lock in share mode;
- 排他锁（X）：select * from table_name where … for update。

#### InnoDB 行锁算法

#### 事务隔离级别

MySQL 的 4 种隔离级别：

- Read uncommitted（读未提交）: 在该隔离级别，所有事务都可以看到其它未提交事务的执行结果。可能会出现脏读。
- Read Committed（读已提交，简称： RC）：一个事务只能看见已经提交事务所做的改变。因为同一事务的其它实例在该实例处理期间可能会有新的 commit，所以可能出现幻读。
- Repeatable Read（可重复读，简称：RR）：这是 MySQL 的默认事务隔离级别，它确保同一事务的多个实例在并发读取数据时，会看到同样的数据行。消除了脏读、不可重复读，默认也不会出现幻读。
- Serializable（串行）：这是最高的隔离级别，它通过强制事务排序，使之不可能相互冲突，从而解决幻读问题。

> 这里解释一下脏读和幻读：
>
> - 脏读：读取未提交的事务。
> - 幻读：一个事务按相同的查询条件重新读取以前检索过的数据，却发现其他事务插入了满足其查询条件的新数据。

#### RC 隔离级别下的行锁实验

##### 通过非索引字段查询

要想分析某条 SQL 是怎么加锁的，如果其他信息都不知道，那就得分几种情况了，不同情况加锁的方式也各不一样，比较常见的一些情况如下：

- RC 隔离级别，a 字段没索引。
- RC 隔离级别，a 字段有唯一索引。
- RC 隔离级别，a 字段有非唯一索引。
- RR 隔离级别，a 字段没索引。
- RR 隔离级别，a 字段有唯一索引。
- RR 隔离级别，a 字段有非唯一索引。

**没有索引的情况下，InnoDB 的当前读会对所有记录都加锁。所以在工作中应该特别注意 InnoDB 这一特性，否则可能会产生大量的锁冲突。**

##### 通过唯一索引查询

**如果查询的条件是唯一索引，那么 SQL 需要在满足条件的唯一索引上加锁，并且会在对应的聚簇索引上加锁。**

##### 通过非唯一索引查询

**如果查询的条件是非唯一索引，那么 SQL 需要在满足条件的非唯一索引上都加上锁，并且会在它们对应的聚簇索引上加锁。**

### 3  答疑篇

#### 关于条件字段做函数操作的问题

下面这条 SQL 走不了索引：

```mysql
select * from t1 where data(c) = '2019-05-21';
```

原因是，索引树中存储的是列的实际值和主键值。如果拿 ‘2019-05-21’ 去匹配，将无法定位到索引树中的值。

有同学就问，如果把语句 “=” 前面的 DATE_FORMAT 按照实际值转换，比如：

```sql
select * from t1 where DATE_FORMAT(c,'%Y-%m-%d %H:%i:%s')='2019-05-21 00:00:00';
```

已经转换成跟实际值一样，为什么也不走索引？

原因是：这条 SQL 执行过程是先找出所有 c 字段的值，然后经过 DATE_FORMAT 加工后形成新的值集合（并且是无序的），因此最终是这个新值集合跟 ‘2019-05-21’ 进行比较。而索引中存储的是 c 字段的集合，所以无法使用索引过滤。

#### 是不是所有的辅助索引都需要回表？

并不是所有的辅助索引查询都需要回表，如果条件字段跟需要查询的字段在一颗索引树上，可以直接在索引树上找到结果，就不用回表了，通常称为覆盖索引。

#### Btree 索引和 HASH 索引有什么区别？

两者区别：

- Btree 索引可能需要多次运用折半查找来找到对应的数据库；
- 而 HASH 索引是通过 HASH 函数，计算出 HASH 值，在表中找出对应的数据。

优缺点对比：

- 大量不同数据等值精确查询，HASH 索引效率通常比 B+TREE 高；
- 但是 HASH 索引不支持模糊查询、排序、范围查询和联合索引中的最左匹配原则，而这些 Btree 索引都支持。

#### 关于隐式转换的一个问题

表 t1 中 a 字段是字符串类型，有索引，但是如下 SQL 走不了索引：

```sql
select * from t1 where a=1000;
```

原因是：MySQL 会把 a 转换成 int 型，再去做判断，也相当于对索引字段做函数操作，导致不能使用索引。

有同学提问：表 t1 中，为 int 型的 b 字段，有索引，下面 SQL 为什么可以走索引？

```sql
select * from t1 where b='1';
```

因为：MySQL 中如果条件字段和后面的值不是同一类型，那么优先将字符串转成数字类型做比较。因此上面的 SQL 是转的 “=” 后面的值，而不是字段 b，因此 SQL 其实等价于：

```sql
select * from t1 where b=1;
```

### 4  间隙锁的意义

### 5  为什么会出现死锁？

**死锁是指两个或者多个事务在同一资源上相互占用，并请求锁定对方占用的资源，从而导致恶性循环的现象。**

InnoDB 中解决死锁问题有两种方式：

1. 检测到死锁的循环依赖，立即返回一个错误（这个报错内容请看下面的实验），将参数 innodb_deadlock_detect 设置为 on 表示开启这个逻辑；
2. 等查询的时间达到锁等待超时的设定后放弃锁请求。这个超时时间由 innodb_lock_wait_timeout 来控制。默认是 50 秒。

#### 同一张表中

不同线程并发访问同一张表的多行数据，未按顺序访问导致死锁。

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| begin;                                                       | begin;                                                       |
| select * from t18 where a=1 for update; … 1 row in set (0.00 sec) | select * from t18 where a=2 for update; … 1 row in set (0.00 sec) |
| select * from t18 where a=2 for update;/* SQL1 */ （等待）   |                                                              |
| **（session2 提示死锁回滚后，SQL1 成功返回结构）**           | select * from t18 where a=1 for update; **ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting transaction** |
| commit;                                                      | commit;                                                      |

####  不同表之间

不同线程并发访问多个表时，未按顺序访问导致死锁：

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| begin;                                                       | begin;                                                       |
| select * from t18 where a=1 for update; … 1 row in set (0.00 sec) | select * from t18_1 where a=1 for update; … 1 row in set (0.00 sec) |
| select * from t18_1 where a=1 for update;/* SQL2 */ 等待     |                                                              |
| **（session2 提示死锁回滚后，SQL1 成功返回结构）**           | select * from t18 where a=1 for update; **ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting transaction** |
| commit;                                                      | commit;                                                      |

#### 如何降低死锁概率？

那么应该怎样降低出现死锁的概率呢？这里总结了如下一些经验：

1. 更新 SQL 的 where 条件尽量用索引；
2. 基于 primary 或 unique key 更新数据；
3. 减少范围更新，尤其非主键、非唯一索引上的范围更新；
4. 加锁顺序一致，尽可能一次性锁定所有需要行；
5. 将 RR 隔离级别调整为 RC 隔离级别。

#### 分析死锁的方法

nnoDB 中，可以使用 SHOW INNODB STATUS 命令来查看最后一个死锁的信息。我们可以尝试用下这个命令获取一些死锁信息，如下：

```sql
show engine innodb status\G
```

## 事务

### 1  数据库忽然断电会丢数据吗？

#### 什么是事务？

事务就是一组原子性的 SQL 查询，或者说一个独立的工作单元。如果数据库引擎能够成功地对数据库应用该组查询的全部语句，那么就执行该组查询。如果其中有任何一条语句因为崩溃或其他原因无法执行，那么所有的语句都不会执行。也就是说，事务内的语句，要么全部执行成功，要么全部执行失败。

ACID 特性：

-  atomicity（原子性） ：要么全执行，要么全都不执行；
-  consistency（一致性）：在事务开始和完成时，数据都必须保持一致状态；
-  isolation（隔离性） ：事务处理过程中的中间状态对外部是不可见的；
-  durability（持久性） ：事务完成之后，它对于数据的修改是永久性的。

#### Redo log

Redo log 称为重做日志，用于记录事务操作变化，记录的是数据被修改之后的值。

Redo log 由两部分组成：

- 内存中的重做日志缓冲（redo log buffer）
- 重做日志文件（redo log file）

每次数据更新会先更新 redo log buffer，然后根据 innodb_flush_log_at_trx_commit 来控制 redo log buffer 更新到 redo log file 的时机。innodb_flush_log_at_trx_commit 有三个值可选：

0：事务提交时，在事务提交时，每秒触发一次 redo log buffer 写磁盘操作，并调用操作系统 fsync 刷新 IO 缓存。

1：事务提交时，InnoDB 立即将缓存中的 redo 日志写到日志文件中，并调用操作系统 fsync 刷新 IO 缓存；

2：事务提交时，InnoDB 立即将缓存中的 redo 日志写到日志文件中，但不是马上调用 fsync 刷新 IO 缓存，而是每秒只做一次磁盘 IO 缓存刷新操作。

innodb_flush_log_at_trx_commit 参数的默认值是 1，也就是每个事务提交的时候都会从 log buffer 写更新记录到日志文件，而且会刷新磁盘缓存，这完全满足事务持久化的要求，是最安全的，但是这样会有比较大的性能损失。

将参数设置为 0 时，如果数据库崩溃，最后 1秒钟的 redo log 可能会由于未及时写入磁盘文件而丢失，这种方式尽管效率最高，但是最不安全。

将参数设置为 2 时，如果数据库崩溃，由于已经执行了重做日志写入磁盘的操作，只是没有做磁盘 IO 刷新操作，因此，只要不发生操作系统奔溃，数据就不会丢失，这种方式是对性能和安全的一种折中处理。

#### Binlog

二进制日志（binlog）记录了所有的 DDL（数据定义语句）和 DML（数据操纵语句），但是不包括 select 和 show 这类操作。Binlog 有以下几个作用：

- 恢复：数据恢复时可以使用二进制日志
- 复制：通过传输二进制日志到从库，然后进行恢复，以实现主从同步
- 审计：可以通过二进制日志进行审计数据的变更操作

可以通过参数 sync_binlog 来控制累积多少个事务后才将二进制日志 fsync 到磁盘。

- sync_binlog=0，表示每次提交事务都只write，不fsync
- sync_binlog=1，表示每次提交事务都会执行fsync
- sync_binlog=N(N>1)，表示每次提交事务都write，累积N个事务后才fsync

比如要加快写入数据的速度或者机器磁盘 IO 瓶颈时，可以将 sync_binlog 设置成大于 1 的值，**但是如果设置为 N(N>1)时，如果数据库崩溃，可能会丢失最近 N 个事务的 binlog。**

### 2  MVCC怎么实现的？

#### 隐藏列

对于 InnoDB ，每行记录除了我们创建的字段外，其实还包含 3 个隐藏的列：

- ROW ID：隐藏的自增 ID，如果表没有主键，InnoDB 会自动以 ROW ID 产生一个聚集索引树。
- 事务 ID：记录最后一次修改该记录的事务 ID。
- 回滚指针：指向这条记录的上一个版本。

#### Undo log

redo log，它记录了事务操作变化。但是事务有时是需要回滚的，这时，Undo log 就发挥了作用。Undo log 是逻辑日志，将数据库逻辑地恢复到原来的样子，所有修改都被逻辑的取消了。

也就是如果是 insert 操作，其对应的回滚操作就是 delete；

如果是 delete，则对应的回滚操作是 insert；

如果是 update，则对应的回滚操作是一个反向的 update 操作。

Undo log 的作用除了回滚操作，Undo log 的另一个作用是 MVCC，InnoDB 存储引擎中 MVCC 的实现是通过 Undo log 来完成的。当用户读取一行记录时，若该记录已经被其它事务占用，当前事务可以通过 Undo log 读取之前的行版本信息，因为没有事务需要对历史的数据进行修改操作，所以也不需要加锁，以此来实现非锁定读取。

#### Read View

ead View 是指事务进行快照读操作的那一刻，产生数据库系统当前活跃事务列表的一个快照。

Read View 中大致包含以下内容：

- trx_ids：数据库系统当前活跃事务 ID 集合；
- low_limit_id：活跃事务中最大的事务 ID +1；
- up_limt_id：活跃事务总最小的事务 ID；
- creator_trx_id：创建这个 Read View 的事务 ID。

比如某个事务，创建了 Read View，那么它的 creator_trx_id 就为这个事务的 ID，假如需要访问某一行，假设这一行记录的隐藏事务 ID 为 t_id，那么可能出现的情况如下：

- 如果 t_id < up_limt_id，说明这行记录在这些活跃的事务创建之前就已经提交了，那么这一行记录对该事务是可见的。
- 如果 t_id >= low_limt_id，说明这行记录在这些活跃的事务开始之后创建的，那么这一行记录对该事物是不可见的。
- 如果 up_limit_id <= t_id < low_limit_id，说明这行记录可能是在这些活跃的事务中创建的，如果 t_id 也同时在 trx_ids 中，则说明 t_id 还未提交，那么这一行记录对该事物是不可见的；如果 t_id 不在 trx_ids 中，则说明事务 t_id 已经提交了，那么这一行记录对该事物是可见的。

对于不可见的记录，都是通过查询 Undo log 来查询老的记录。

了解了上面的原理，我们知道了，Read View 规则帮我们判断当前版本的数据是否可见。下面，我们分析下当查询一条记录时，大致的步骤：

- 获取事务本身的事务 ID；
- 获取 Read View；
- 查询得到的数据，然后与 Read View 中的事务版本号进行比较；
- 如果能查询，则直接查询对应的记录；如果不能直接查询，则通过 Undo Log 中获取历史快照；
- 最终返回结果。

另外需要补充的一点就是，在 RR 和 RC 隔离级别下，获取 Read View 的时机也是不一样的：

在可重复读隔离级别（RR）下，同一个事务中，查询语句只是在第一个读请求发起时获取 Read View，而后面相同的查询语句都会使用这个 Read View。

#### 什么是MVCC？

| session1 | session2                                                     |                                                              |
| :------- | :----------------------------------------------------------- | ------------------------------------------------------------ |
| 1        | set session transaction_isolation=‘READ-COMMITTED’;/* 设置会话隔离级别为 RC*/ | set session transaction_isolation=‘READ-COMMITTED’;/* 设置会话隔离级别为 RC*/ |
| 2        | select * from t20;![图片描述](https://img1.sycdn.imooc.com/5d832c8900016c3c03340181.png) |                                                              |
| 3        | begin;                                                       |                                                              |
| 4        | update t20 set b=666 where a=1;                              |                                                              |
| 5        |                                                              | begin;                                                       |
| 6        |                                                              | select * from t20;![图片描述](https://img1.sycdn.imooc.com/5d832caa00016c3c03340181.png) |
| 7        | commit;                                                      |                                                              |
| 8        |                                                              | select * from t20;![图片描述](https://img1.sycdn.imooc.com/5d832cc7000103f903520175.png) |
| 9        |                                                              | commit;                                                      |

在 session1 更新了 a=1 这行记录，但还没提交的情况下，在 session2 中，满足 a=1 这条记录，b 的值还是原始值 1，而不是 session 1 更新之后的 666，那么在数据库层面，这是怎么实现的呢？

其实 InnoDB 就是通过 MVCC 和 Undo log 来实现的。

什么是 MVCC 呢？

**MVCC， 即多版本并发控制。MVCC 的实现，是通过保存数据在某个时间点的快照来实现的，也就是说，不管需要执行多长时间，每个事务看到的数据都是一致的。根据事务开始的时间不同，每个事务对同一张表，同一时刻看到的数据可能是不一样的。**

也就是上面实验第 6 步中，为什么 session2 查询的结果还是 session1 修改之前的记录。

#### MVCC的实现原理

#### MVCC的优势

MVCC 最大的好处是读不加锁，读写不冲突，极大地增加了 MySQL 的并发性。

通过 MVCC，保证了事务 ACID 中的 I（隔离性）特性。

#### 总结

MVCC 实现的原理大致是：

InnoDB 每一行数据都有一个隐藏的回滚指针，用于指向该行修改前的最后一个历史版本，这个历史版本存放在 Undo log 中。如果要执行更新操作，会将原记录放入 Undo log 中，并通过隐藏的回滚指针指向 Undo log 中的原记录。其它事务此时需要查询时，就是查询 Undo log 中这行数据的最后一个历史版本。

### 3  养成好的事务习惯

#### 不好的事物习惯

##### 在循环中提交

##### 不关注同一个事务里语句顺序

比如 A 在超市购买 100 元的商品，付款操作可以简化为：

| 序号 | 操作                     |
| :--- | :----------------------- |
| 1    | A 的账户中扣除 100       |
| 2    | 超市的账户增加 100       |
| 3    | 在超市系统中记录一条日志 |

很多时候我们会按上面的 SQL 步骤放入一个事务里执行，不关注里面语句的顺序。实际可以优化的。

根据两阶段锁，整个事务里面涉及的锁，需要等到事务提交时才会释放。因此我们在**同一个事务中，可以把没锁或者锁范围小的语句放在事务前面执行，而锁定范围大的语句放在后面执行。**

##### 不关注不同事务访问资源的顺序

几种产生死锁的原因，其中有两条就跟不同事务访问资源顺序有关，我们来回顾一下：

- 不同线程并发访问同一张表的多行数据，未按顺序访问导致死锁。
- 不同线程并发访问多个表时，未按顺序访问导致死锁。

如果不关注并发访问的不同事务中访问资源的顺序，就会增大出现死锁的概率。

因此，为了降低死锁，我们需要去关注不同事务访问资源的顺序。

##### 不关注事务隔离级别

如果完全不关注自己业务使用的 MySQL 是什么隔离级别，可能会降低程序的并发能力或者导致死锁。

比如业务场景完全能接受幻读，如果要求更高的 QPS，使用 RR 隔离级别显然不是最好的选择，因此可以改为 RC 隔离级别。

而如果业务使用的是 RR 隔离级别，可能由于间隙锁导致死锁（可参考第 18 节 2.3 中的例子），因此也应该在程序编写时关注 RR 隔离级别下是否会有间隙锁。

因此，为了更高的并发和降低死锁概率，在创建事务前，也应该去关注自己业务的数据库是什么事务隔离级别。

##### 在事务中混合使用存储引擎

在事务中混合使用事务型（比如 InnoDB）和非事务型（比如 MyISAM）表，如果是正常提交，到没什么问题。

但是，如果该事务回滚了，事务型的表可以正常回滚，而非事务型的表的变更就无法回滚了。这种情况就会导致数据不正常，并且事务最终的结果也难以确定。

### 4  细聊分布式事务

#### 认识分布式事务

分布式事务是指一个大的事务由很多小操作组成，小操作分布在不同的服务器上或者不同的应用程序上。分布式事务需要保证这些小操作要么全部成功，要么全部失败。MySQL 从 5.0.3 开始支持分布式事务。

分布式事务使用两阶段提交协议：

- 第一阶段：所有分支事务都开始准备，告诉事务管理器自己已经准备好了；
- 第二阶段：确定是 rollback 还是 commit，如果有一个节点不能提交，则所有节点都要回滚。

与本地事务不同点在于：分布式事务需要多一次 prepare 操作，等收到所有节点的确定信息后，再进行 commit 或者 rollback。

#### MySQL自带的分布式事务

#### 结合中间件实现分布式事务

具体实现方式可以拿上面网上购书的例子来说：

订单业务程序处理完增加订单的操作后，将减库存操作发送到消息队列中间件中（比如：Rocketmq），订单业务程序完成提交。然后库存业务程序检查到消息队列有减对应商品库存的信息，就开始执行减库存操作。库存业务执行完减库存操作，再发送一条消息给消息队列中间件：内容是已经减掉库存。

## MySQL的一些其他经验

### 1  如何预防SQL注入？

#### 认识SQL

SQL 注入是利用某些数据库的外部接口将用户数据插入到实际的数据库操作语句中，从而达到入侵数据库乃至操作系统的目的。

SQL 注入产生的主要原因是：**程序对用户输入的数据没有进行严格的过滤，导致非法数据库查询语句的执行**。

SQL 注入具有很大的危害，可能会导致攻击者非法入侵系统，或者盗取数据，甚至清空数据等。

#### 如何进行 SQL 注入攻击？

但是可以通过 SQL 注入的方式，让攻击者不知道密码的情况也可以通过验证，如下：

把 URL 改成：

```http
http://172.16.4.15/index.php?username=mt' or'1=1
```

#### 如何预防SQL注入？

##### 控制变输入变量的格式

##### 转义特殊字符

### 2  主键是否需要设计为自增？

#### 关于自增主键

发现表结构中包含 AUTO_INCREMENT=2，在下一次执行未指定 id 字段的 insert 语句时， 主键 id 会自动获取到这个值。

#### 主键和聚集索引的关系

它们的区别：

- 在 InnoDB 中，聚集索引不一定是主键，但是主键一定是聚集索引：原因是如果没有定义主键，聚集索引可能是第一个不允许为 null 的唯一索引，如果也没有这样的唯一索引，InnoDB 会选择内置 6 字节长的 ROWID 作为隐含的聚集索引。
- InnoDB 的数据是按照主键顺序存放的，而聚集索引就是按照每张表的主键构造一颗 B+ 树，它的叶子节点存放的是整行数据。
- 每张 InnoDB 表都有一个聚集索引，但是不一定有主键。

#### 主键是否需要设置为自增

**因此如果业务输入的主键都是随机数字，那么写入数据时很可能会导致数据页频繁分裂，从而影响写入效率。**

**而如果设置主键是自增，那么每一次都是在聚集索引的最后增加，当一页写满，就会自动开辟一个新页，不会有聚集索引树分裂这一步，效率会比随机主键高很多。这也是很多建表规范要求主键自增的原因。**

除了要求主键自增外，最好主键也要无业务意义，原因是防止主键更新而导致页分裂的情况。

### 3  MySQL是否需要开启查询缓存？

从 MySQL 4.1 开始，增加了查询缓存（Query Cache，简称 QC）的功能，它会存储 select 语句的文本以及发送到客户端的结果。如果下一次收到一个相同的查询，就会从查询缓存中获得查询结果。

#### 认识QC

QC 需要缓存最新数据结果，因此表数据发生任何变化（insert、update、delete 等操作时），都会导致 QC 被刷新。

查询缓存相关的参数主要有：

```sql
show variables like "%query_cache%";
```

**参数**

have_query_cache：服务器在安装时是否已经配置了高速缓存
query_cache_limit：单条查询能够使用的缓存区大小
query_cache_min_res_unit：查询缓存分配内存块的最小单位
query_cache_size：缓存区的大小，单位为 MB
query_cache_type：缓存类型，有三个值可选：

- 0 或者 off：关闭缓存
- 1 或者 on：打开缓存
- 2 或者 demand：只缓存带有 sql_cache 的 select 语句。

query_cache_wlock_invalidate：如果某个数据表被其它的连接锁住，是否仍然从查询缓存中返回结果

通过下面命令，可以监视查询缓存的使用情况：

```sql
show global status like "qcache%";
```

**参数**

Qcache_free_blocks：查询缓存中的空闲内存块的数目

Qcache_free_memory：查询缓存的空闲内存总数

Qcache_hits：缓存命中次数

Qcache_inserts：被加入到缓存中的查询数目

Qcache_lowmem_prunes：因为缺少内存而从缓存中删除的查询数目

Qcache_not_cached：没有被缓存的查询数目

Qcache_queries_in_cache：在缓存中已注册的查询数目

Qcache_total_blocks：查询缓存中的块的总数目

#### QC 的优劣

##### QC优势

提高查询速度：使用查询缓存在单行数据的表中搜索要比不使用查询缓存快 238%

##### QC劣势

- 如执行的 SQL 都很简单（比如从只有一行的表中查询数据），但每次查询都不一样的话，打开 QC 后，额外的开销为 13% 左右；
- 如果表数据发生了修改，使用该表的所有缓存查询都将失效，并且从缓存中删除；
- QC 要求前后两次请求的 SQL 完全一样，不同数据库、不同协议版本或不同默认字符集的查询，都会被认为是不同的查询。甚至包括大小写，比如下面两条 SQL ，查询缓存就会认为是两个不同的查询：

```sql
SELECT * FROM tbl_name
Select * from tbl_name
```

- 每次更新 QC 的内存块都需要进行锁定；
- 可能会导致 SQL 查询时间不稳定。

#### 是否需要开启QC

如果线上环境中 99% 以上都是只读，很少更新，可以考虑全局开启 QC，也就是设置 query_cache_type 为 1。

很多时候，我们希望缓存的是几张更新频率很低的表，其它表不考虑使用查询缓存，就可以考虑将 query_cache_type 设置成 2 或者 DEMAND，这样就只缓存下面这类 SQL：

```sql
select sql_cache ......;
```

结合前面讲的 QC 的劣势，其它情况就不建议开启 QC 了。

#### 怎样开启和关闭QC

全局开启 QC：

在配置文件 my.cnf 中设置：

```sql
query_cache_type = 1
query_cache_size = 50M
```

只开启部分表的 QC：

在配置文件 my.cnf 中设置：

```sql
query_cache_type = 2
query_cache_size = 50M
```

前面也提到了，这种情况，如果要使用 QC，查询语句需要加上 select sql_cache。

关闭 QC：

在配置文件 my.cnf 中设置：

```sql
query_cache_type = 0
query_cache_size = 0
```

或者源码编译安装 MySQL 的话，编译时增加参数 --without-query-cache 即可。

#### 开启QC的注意事项

如果要开启 QC，建议不要设置过大，通常几十兆就好。如果设置过大，会增加维护缓存所需要的开销。

另外要注意一些即使开启 QC 也不能使用 QC 的场景）：

- 分区表不支持，如果涉及分区表的查询，将自动禁用查询缓存
- 子查询或者外层查询
- 存储过程、触发器中使用的 SQL
- 读取系统库时
- 类似下面 SQL 时：

```sql
select ... lock in share mode
select ... for update
```

- 用到临时表
- 产生了 warning 的查询
- 显示增加了 SQL_NO_CACHE 关键字的
- 如果没有全部库、表的 select 权限，则也不会使用 QC
- 使用了一些函数：比如 now ()，user ()，password () 等

了解上面的场景，我们就能知道：开了查询缓存，前后 SQL 一模一样，为什么后面这一次执行也使用不了缓存的原因了。

### 4  使用读写分离需要注意哪些？

通常我们说的 MySQL 读写分离是指：对于修改操作在主库上执行，而对于查询操作，在从库上执行。主要目的是分担主库的压力。

但是读写分离有时也会存在问题，比如：主从延迟时，读取的从库数据不是最新的，对应的业务场景比如：

你网购的一个商品，付完款之后，因为主从延迟，第一时间还查询不到订单（查询的从库），即使等一段时间能看到订单，但是相信这种情况很多用户是不能接受的。

本节就一起来讨论一下：读写分离需要注意哪些问题。

通常情况下，读写分离都是依赖主从复制，因此，我们先来看看主从复制的原理，也能方便我们理解为什么会出现主从延迟的现象。

#### 主从复制的原理

##### MySQL异步复制

传统的 MySQL 主从复制是异步的，因此也称为异步复制，MySQL 异步复制的原理如下：

- 在主库开启 binlog 的情况下
- 如果主库有增删改的语句，会记录到 binlog 中
- 主库通过 IO 线程把 binlog 里面的内容传给从库的中继日志（relay log）中
- 主库给客户端返回 commit 成功（这里不会管从库是否已经收到了事务的 binlog）
- 从库的 SQL 线程负责读取它的 relay log 里的信息并应用到从库数据库中

**在主库上并行运行的更新 SQL，由于从库只有单个 SQL 线程去消化 relay log，因此更新的 SQL 在从库只能串行执行。这也是很多情况下，会出现主从延迟的原因。**

当然，从 5.6 开始，MySQL 支持了每个库可以配置单独的 SQL 线程来消化 relay log，在 5.7 又增加了基于组提交的并行复制，大大改善了主从延迟的问题。

##### MySQL半同步复制

在 MySQL 异步复制的基础上，又出现了一种改进的复制方式，称为：半同步复制。其原理如下：

- 在主库开启 binlog 的情况下
- 如果主库有增删改的语句，会记录到 binlog 中
- 主库通过 IO 线程把 binlog 里面的内容传给从库的中继日志（relay log）中
- **从库收到 binlog 后，发送给主库一个 ACK，表示收到了**
- **主库收到这个 ACK 以后，才能给客户端返回 commit 成功**
- 从库的 SQL 线程负责读取它的 relay log 里的信息并应用到从库数据库中

跟传统的异步复制相比，**半同步复制保证了所有给客户端发送过确认提交的事务**，从库都已经收到这个日志了。

#### 常见的读写分离方式

##### 通过程序

开发通过配置程序来决定修改操作走主库，查询操作走从库。这种方式直连数据库，优点是性能会好点，缺点是配置麻烦。

但是需要注意的是：从库需要设置为 read_only，防止配置错误在从库写入了数据。

##### 通过中间件

通过中间件实现读写分离，目前算是一种主流的方式。拿 MyCAT 举例：

在 schema.xml 文件中，dataHost 标签 balance 属性的值，决定了是否启用读写分离。

balance 各个值及对应的读写方法如下：

- 0：不开启读写分离，读操作发送到 writehost
- 1：全部的 readhost 与 stand by writehost 参与 select 语句的负载均衡
- 2：所有读操作都随机在 writehost、readhost上分发
- 3：所有读请求随机分发到 writerhost 对应的 readhost 执行，writehost 不负担读压力

因此可以根据实际情况选择上面合适的读写分离策略。

##### 什么情况下会出现主从延迟？

对于读写分离场景，最大的问题就是：主从延迟。那么在哪些情况下会出现主从延迟呢？这里大致总结一下可能导致主从延迟的场景：

- 大表 DDL
- 大事务
- 主库 DML 并发大
- 从库配置差
- 表上无主键

#### 读写分离怎样应对主从延迟？

##### 判断主从是否延迟

有些业务场景，如果所有请求都落在主库，主库压力会很大，但是在读写分离的情况，又不希望主从存在延迟的时候去读取从库。这种情况，就可以考虑查询时，先判断主从是否存在延迟，如果存在延迟，则查询落在主库，如果没延迟，则查询语句落在从库。

第一种方法：判断 Seconds_Behind_Master 是否等于 0。

如果 Seconds_Behind_Master =0，则查询从库，如果大于 0，则查询主库。

##### 采用半同步复制

跟传统的异步复制相比，半同步复制保证了所有给客户端发送过确认提交的事务，从库都已经收到这个日志了。因此出现延迟的概率会小很多，当然实际生产应用时，建议结合上面讲的位点或 GTID 判断。

##### 等待同步完成

有人可能会觉得：这种方式谁会用啊？实际可以应用在内部人员看的报表业务上。因为报表可能涉及的 SQL 都比较复杂，存在延迟就考虑去查询主库，可能会对其它线上业务有影响，因此可以等待从库同步完成，再查询从库。

### 5  哪些情况需要考虑读写分库分表？

MySQL 分库分表是指：把 MySQL 数据库物理地拆分到多个实例或者机器上去。从而降低单台 MySQL 实例的负载。那么这里讲的拆分究竟以什么方式实现的呢？这里就先来看看 MySQL 分库分表的拆分方法：

#### MySQL分库分表拆分方法

MySQL 分库分表拆分方法分为：垂直拆分和水平拆分。

##### 垂直拆分

垂直拆分一般是指下面几种情况：

- 有多个业务，每个业务单独分到一个实例里面。
- 在一个实例中有多个库，把这些库分别放到单独的实例中。
- 在一个库中存在过多的表，把这些表拆分到多个库中。
- 把字段过多的表拆分成多个表，每张表包含一部分字段。

##### 水平拆分

所谓水平拆分，就是把同一张表分为多张表结构相同的表，每张表里存储一部分数据。而拆分的算法也比较多，常见的就是取模、范围、和全局表等。

#### 哪些情况需要考虑分库分表？

##### 数据量过大，影响了运维操作

如果数据量比较大，其对数据库影响也会比较大，最常见的比如：影响备份、大表 DDL 导致主从长时间延迟等。下面仔细讲讲这两种情况：

备份：如果单张表或者单个实例数据量太大，那备份可能需要占用大量的 IO 和磁盘空间，并且持续时间还会比较久。曾经听说过有公司的单个实例的备份从凌晨持续到当天下午的情况，这种场景下，尽管在从库备份，如果开启了读写分离，对业务的影响也是比较大的。

DDL 导致主从长时间延迟：大表执行 DDL 不但会产生 MDL 写锁，并且还会导致主从延迟。

##### 把修改频繁的字段拆分出来

比如电商业务场景的用户表，这张表可能包含了用户唯一标识 ID，用户名，昵称，联系方式，性别，出生日期，注册时间，积分等等。这些字段中，其实也只有积分会变更的比较频繁，因此可以把积分字段独立出来，然后加上用户唯一标识 ID。这样，用户表的更新次数就大大降低了。

##### 把大字段拆分出去

比如电商业务的商品表，表里可能包含了商品的价格，生产日期，产地，供应商，商品详情和使用说明等。我们会发现商品详情和使用说明的字段特别大，可能字段类型需要配置成 text 或 blob，类似的字段占用的磁盘空间也是比较大。维护起来会比较麻烦，因此这种情况下，建议把这类大字段水平拆分出来。

##### 增长比较快的

还是拿电商业务举例，比如订单表，如果预估未来增长速度会比较快，那么可以考虑提前对订单表进行分库分表。防止单表增长过快。

##### 降低不同库或者表的相互影响

电商业务会涉及多个模块，比如会员，订单，库存等。比如搞了一个活动，扫码送会员积分，如果会员相关的表和订单相关的表在一个库中，那么很可能会员的活动会影响到订单业务。因此可以考虑将这些模块对应的表拆分到不同的库中。避免不同业务表或者库的相互影响。

#### 分库分表的实现

##### 通过程序

如我已经通过垂直拆分，将不同业务表放在不同的业务库中。这种情况，程序只要每个业务配置不同的 database 即可。

如果是水平拆分，业务可以通过程序实现，比如按照用户名首字母分，讲不同首字母的名字过来的请求放在不同的库中。

##### 通过数据库中间件

目前比较常见的就是通过数据库中间件，比如 MyCAT。这种对业务程序来说就比较友好了，某个业务库，程序端只要连逻辑库，而后端是通过 MyCAT，将这一个逻辑库指向多个物理库。然后通过分片字段，决定数据应该放在哪个分片里。

### 6  如何安全高效地删除大量无用数据？

为了防止单张表过大，很多情况下，我们需要对历史数据做归档或者删除操作。那么如何高效删除大量无用数据呢？

#### 共享表空间和独立表空间

InnoDB 数据是按照表空间进行存放的，其表空间分为共享表空间和独立表空间。

##### 共享表空间

##### 独立表空间

#### 几种数据删除形式

##### 删除表

如果是某张表（假如表名为 t29）的数据和表结构都不需要使用了，那么可以考虑 drop 掉。出于安全考虑，建议的步骤如下：

- 1、首先将表名改为 t29_bak_20191011：

```sql
alter table t29 rename t29_bak_20191011;
```

- 2、然后等待半个月，观察是否有程序因为找不到表 t29 而报错；
- 3、如果没有跟表 t29 相关的报错，则半个月后直接 drop 掉 t29_bak_20191011：

```sql
drop table t29_bak_20191011;
```

##### 清空表

如果是某张表（假如表名为 t29）的历史数据不需要使用了，要做一次清空，则可以考虑使用 truncate。

建议的步骤如下：

1、创建一张与 t29 表结构相同的临时表：

```sql
create table t29_bak_20191011 like t29;
```

2、并将数据拷贝到临时表：

```sql
insert into t29_bak_20191011  select * from t29;
```

3、再清空该表：

```sql
truncate table t29;
```

4、如果空间不够，观察半个月后，考虑转移 t29_bak_20191011 的数据到备份机器上。然后删除表 t29_bak_20191011：

```sql
drop table t29_bak_20191011;
```

我不止一次遇到类似的情况：需要清空表而使用 delete from table_name，导致主从延迟和磁盘 IO 跑满的情况。 原因是 binlog 为行模式的情况下，执行全表 delete 会生成每一行对应的删除操作，因此可能导致单个删除事务非常大。而 truncate 可以理解为 drop + create，在 binlog 为 row 模式的情况下，也只会产生一行 truncate 操作。所以，建议清空表时使用 truncate 而不使用 delete。

##### 非分区表删除部分记录

##### 分区表删除部分分区

### 7  使用MySQL时，应用层可以这么优化

#### 使用连接池

MySQL 如果频繁创建和断开连接，那 MySQL 的开销会比较大，可能会占用过多的服务器内存资源，甚至导致响应时间变慢。此时就可以考虑使用连接池来改进性能。

连接池可以理解为：创建一些持久连接的“池”，新的请求可以使用这些连接池，减少创建和断开连接的次数。

其大致原理是：

1. 当进程启动时，创建相应的数据库连接池对象；
2. 如果程序需要请求数据库，则直接从连接池获取到一个连接；
3. 数据库请求完成后，释放数据库连接池。

**那么连接池会不会导致服务器连接过多呢？**

通常情况下，连接池不会导致服务器连接过多，因为它们会在进程间排队和共享连接。

#### 减少对MySQL的访问

避免对同一行数据做重复检索，比如查询某个用户信息。

首先查出这个用户的联系方式：

```sql
select phone from user_info where user_id=111;
```

然后再查出这个用户的姓名：

```sql
select name from user_info where user_id=111;
```

显然上面的方式并不是最优的，可以将两条 SQL 合并成一条：

```sql
select phone,name  from user_info where user_id=111;
预览
```

然后返回给客户端。这样跟数据库建立连接的次数从 2 次降低到 1 次，从而节省了部分建立连接所花费的内存和时间。

#### 增加redis缓存层

##### 计数器

使用 Redis 做计数器的场景。当统计数增加，则在 Redis 中执行下面的命令让计数器加 1：

```shell
INCR t1_count
预览
```

通过这种方式缓解在 MySQL 中执行 update 的压力。

##### K-V数据缓存

在 MySQL 中，如果某个字段会被频繁查询，而该字段内容变化的概率又不是很大，就可以考虑使用 Redis 缓存。比如电商业务，查看上个季度某类型商品的销量排行，如果这一个功能放在主页，点击量可能会非常高，因此可以考虑放在 Redis 中。

##### 消息队列

如果用户在 APP 上点击注册用户名密码，需要填下手机号验证，程序会将该用户的手机号放在 Redis 的 list 中，然后另外一个程序一直去消费 list 中的手机号，取出手机后，则调用第三方短信接口，发送手机短信息给到用户。

#### 单表过大及时归档

比如单张表过大，可能有下面这些影响：

1. 在修改表结构时导致长时间主从延迟；
2. 备份时间过久；
3. 查询速度可能也会变慢。

因此，可以考虑对历史数据归档（比如日志数据），控制单表的数据量。

#### 代码层读写分离

在配置了 MySQL 主从环境的情况下，可以考虑使用读写分离，通过程序配置的这种方式。更新走主库，查询走从库。

#### 表的索引提前规划

### 8  MySQL整体优化思路

#### 硬件相关优化

##### CPU相关

##### 内存相关

##### 磁盘相关

#### 系统层面优化

##### 调整I/O调度算法

##### 文件系统选择

##### 调整内核参数

#### MySQL层面优化

##### 参数优化

##### 设计优化

1. 使用 InnoDB 存储引擎，不建议使用 MyISAM 存储引擎；
2. 预估表数据量和访问量，如果数据量或者访问量比较大，则需要提前考虑分库分表；
3. 指定合适的数据库规范，在设计表、执行 SQL 语句时按照数据库规范来进行。

### 9  MySQL操作规范

#### 命名规范

1.1、表名建议使用有业务意义的英文词汇，必要时可加数字和下划线，并以英文字母开头；

1.2、库、表、字段全部采用小写；

MySQL 在 Linux 下默认是区分大小写的，而在 Windows 下不区分大小写。因此，防止出现问题，建议都设置为小写。

1.3、避免用 MySQL 的保留字， MySQL 保留字请参考[官方手册：9.3 Keywords and Reserved Words](https://dev.mysql.com/doc/refman/5.7/en/keywords.html)；

1.4、命名（包括表名、列名）禁止超过 30 个字符；

1.5、临时库、表名必须以 tmp 为前缀，并以日期为后缀，如：tmp_shop_info_20190404；

1.6、备份库、表必须以 bak 为前缀，并以日期为后缀，如：bak_shop_info_20190404；

1.7、索引命名：

- 非唯一索引必须按照“idx_字段名称”进行命名；
- 唯一索引必须按照“uniq_字段名称”进行命名。

#### 设计规范

2.1、主键：

- 表必须有主键；
- 不使用更新频繁的列做主键；
- 尽量不选择字符串列做主键；
- 不使用 UUID MD5 HASH 做主键；
- 默认使用非空的唯一键。

2.2、如无特殊要求，建议都使用 InnoDB 引擎；

2.3、默认使用 utf8mb4 字符集，数据排序规则使用 utf8mb4_general_ci；

原因：utf8mb4 为万国码，无乱码风险；与 utf8 编码相比，utf8mb4 能支持 Emoji 表情。

2.4、所有表、字段都需要增加 comment 来描述此表、字段所表示的含义；

比如：data_status TINYINT NOT NULL DEFAULT ‘1’ COMMENT ‘1代表记录有效，0代表记录无效’。

2.5、如无说明，表必须包含 create_time 和 update_time 字段，即表必须包含记录创建时间和修改时间的字段；

2.6、用尽量少的存储空间来存数一个字段的数据：

- 能用 int 的就不用 char 或者 varchar；
- 能用 tinyint 的就不用 int；
- 使用 UNSIGNED 存储非负数值；
- 只存储年使用 YEAR 类型；
- 只存储日期使用 DATE 类型。

2.7、存储精确浮点数必须使用 DECIMAL 替代 FLOAT 和 DOUBLE；

原因：在存储的时候，FLOAT 和 DOUBLE 都存在精度损失的问题，很可能在比较值的时候，得到不正确的结果。

2.8、尽可能不使用 TEXT、BLOB 类型；

原因：会浪费更多的磁盘和内存空间，非必要的大量大字段查询会淘汰掉热数据，导致内存命中率急剧降低，影响数据库性能。如果实在有某个字段过长需要使用 TEXT、BLOB 类型，则建议独立出来一张表，用主键来对应，避免影响原表的查询效率。

2.9、禁止在数据库中存储明文密码；

2.10、索引设计规范：

a、需要添加索引的字段

- UPDATE、DELETE 语句的 WHERE 条件列；
- ORDER BY、GROUP BY、DISTINCT 的字段（原因可复习[第 6 节](https://www.imooc.com/read/43/article/686)）；
- 多表 JOIN 的字段（原因可复习[第 8 节](https://www.imooc.com/read/43/article/688)）。

b、单表索引建议控制在 5 个以内；

c、适当配置联合索引；

比如方便查询能走覆盖索引，或者几个字段同时作为条件的概率很高时，当然还有其他很多种情况可以设置联合索引，具体可以复习[第 13 节](https://www.imooc.com/read/43/article/693)。

d、业务上具有唯一性的字段，添加成唯一索引；

遇到过几次字段在业务场景上要求唯一，但是该字段在数据库里的数据却出现了重复。因此在代码层考虑外，还需要在 MySQL 上的对应字段添加唯一索引。

e、在 varchar 字段上建立索引时，建议根据实际文本区分度指定索引长度；

原因：可以降低索引所占用的空间，并且很多时候，比如字符串基本是长度大于 20，但是只要建立长度为 20 的索引，就已经有很高的区分度了。可以使用 count(distinct left(列名, 索引长度))/count(*) 的区分度来确定。

f、索引禁忌：

- 不在低基数列上建立索引，例如：性别字段。
- 不在索引列进行数学运算和函数运算（原因，做函数操作可能会导致使用不了索引，具体可以复习[第 4 节](https://www.imooc.com/read/43/article/684)）

2.11、不建议使用外键；

原因：外键会导致表与表之间耦合，update 与 delete 操作都会涉及相关联的表，十分影响 sql 的性能，甚至会造成死锁。高并发情况下容易造成数据库性能。

2.12、禁止使用存储过程、视图、触发器、Event ；

原因：高并发的情况下，这些功能很可能将数据库拖死，业务逻辑放到服务层具备更好的扩展性。

2.13、单表列数目建议小于 30；

#### SQL语句规范

3.1、避免隐式转换；

3.2、尽量不使用select *,只 select 需要的字段 ；

原因：读取不需要的列会增加 CPU、IO、NET 消耗，并且不能有效的利用覆盖索引。使用 SELECT * 容易在增加或者删除字段后导致程序报错。

3.3、禁止使用 INSERT INTO t_xxx VALUES (xxx)，必须显示指定插入的列属性 ；

原因：容易在增加或者删除字段后导致程序报错。

3.4、尽量不使用负向查询；

比如 not in/like。

3.5、禁止以 % 开头的模糊查询。

原因：使用不了索引。

3.6、禁止单条 SQL 语句同时更新多个表；

3.7、统计记录数使用 select count(*)，而不是 select count(primary_key)或者 select count(普通字段名)；

原因：可能会导致走的索引不是最优的或者导致统计数字不准确。

3.8、建议将子查询转换为关联查询；

3.9、建议应用程序捕获 SQL 异常，并有相应处理；

3.10、SQL 中不建议使用 sleep()，如特殊需求需要用到 sleep()，请提前告知 DBA；

3.11、避免大表的 join。

#### 行为规范

4.1、批量导入、导出数据必须提前通知 DBA 协助观察；

4.2、有可能导致 MySQL QPS 上升的活动，提前告知DBA；

4.3、同一张表的多个 alter 合成一次操作；

4.4、不在业务高峰期批量更新、查询数据库；

4.5、删除表或者库要求尽量先 rename，观察几天，确定对业务没影响，再 drop。在 Extra 字段发现 “Select tables optimized away” 关键字，表示是从 MyISAM 引擎维护的准确行数上获取到的统计值。

而 InnoDB 并不会保留表中的行数，**因为并发事务可能同时读取到不同的行数**。所以执行 count(*) 时都是临时去计算的，会比 MyISAM 引擎慢很多。

对比 MyISAM 引擎和 InnoDB 引擎 count(*) 的区别：

- MyISAM 会维护表的总行数，放在磁盘中，如果有 count(*) 的需求，直接返回这个数据
- 但是 InnoDB 就会去遍历普通索引树，计算表数据总量

##### MySQL 5.7.18 前后 count(*) 的区别

在 MySQL 5.7.18 之前，InnoDB 通过扫描聚簇索引来处理 count(*) 语句。

从 MySQL 5.7.18 开始，通过遍历最小的可用二级索引来处理 count(*) 语句。如果不存在二级索引，则扫描聚簇索引。但是，如果索引记录不完全在缓存池中的话，处理 count(*) 也是比较久的。

新版本为什么会使用二级索引来处理 count(*) 语句呢？

原因是 InnoDB 二级索引树的叶子节点上存放的是主键，而主键索引树的叶子节点上存放的是整行数据，所以二级索引树比主键索引树小。因此优化器基于成本的考虑，优先选择的是二级索引。所以 count(主键) 其实没 count (*) 快。

##### count(1) 比 count(*) 快吗？

count(*) 无论是否包含空值，所有结果都会统计。

而 count(1)中的 1 是恒真表达式，因此也会统计所有结果。

所以 count(1) 和 count(*) 统计结果没差别。

#### 哪些方法可以加快 count()？

##### show table status

```sql
show table status like 't1';
```

##### 用 Redis 做计数器

##### 增加计数表

## MYSQL索引

### 1  为什么添加索引能提高查询速度？

#### 跟索引相关的一些算法

对于 MySQL 而言，使用最频繁的就是 B+ 树索引，所以我们必须要知道 B+ 树的结构，而 B+ 树是借鉴了二分查找法、二叉查找树、平衡二叉树、B 树的一些思想构建的。

##### 二分查找法

二分查找法的查找过程是：将记录按顺序排列，查找时先以有序列的中点位置为比较对象，如果要找的元素值小于该中点元素，则将查询范围缩小为左半部分；如果要找的元素值大于该中点元素，则将查询范围缩小为右半部分。以此类推，直到查到需要的值。

##### 二叉查找树

二分查找法的查找过程是：将记录按顺序排列，查找时先以有序列的中点位置为比较对象，如果要找的元素值小于该中点元素，则将查询范围缩小为左半部分；如果要找的元素值大于该中点元素，则将查询范围缩小为右半部分。以此类推，直到查到需要的值。

##### 平衡二叉树

我们一起看下平衡二叉树的定义：满足二叉查找树的定义，另外必须满足任何节点的两个子树的高度差最大为 1。

##### B 树

B 树可以理解为一个节点可以拥有多于 2 个子节点的多叉查找树。

B 树中同一键值不会出现多次，要么在叶子节点，要么在内节点上。

##### B+ 树

B+ 树是 B 树的变体，定义基本与 B 树一致，与 B 树的不同点：

- 所有叶子节点中包含了全部关键字的信息
- 各叶子节点用指针进行连接
- 非叶子节点上只存储 key 的信息，这样相对 B 树，可以增加每一页中存储 key 的数量。
- B 树是纵向扩展，最终变成一个 “瘦高个”，而 B+ 树是横向扩展的

B+树与B树的结构最大的区别就是：

它的键一定会出现在叶子节点上，同时也有可能在非叶子节点中重复出现。而 B 树中同一键值不会出现多次。

#### B+ 树索引

B+ 树索引就是基于本节前面介绍的 B+ 树发展而来的。在数据库中，B+ 树的高度一般都在 2 ~ 4 层，所以**查找某一行数据最多只需要 2 到 4 次 IO。而没索引的情况，需要逐行扫描，明显效率低很多，这也就是为什么添加索引能提高查询速度。**

B+ 树索引并不能找到一个给定键值的具体行，B+ 树索引能找到的只是被查找数据行所在的页。然后数据库通过把页读入到缓冲池（buffer pool）中，在内存中通过二分查找法进行查找，得到需要的数据。

InnoDB 中 B+ 树索引分为聚集索引和辅助索引，我们再继续了解这两种索引的特点。

#####  聚集索引

InnoDB 的数据是按照主键顺序存放的，而聚集索引就是按照每张表的主键构造一颗 B+ 树，它的叶子节点存放的是整行数据。

InnoDB 的主键一定是聚集索引。如果没有定义主键，聚集索引可能是第一个不允许为 null 的唯一索引，也有可能是 row id。

由于实际的数据页只能按照一颗 B+ 树进行排序，因此每张表只能有一个聚集索引（TokuDB 引擎除外）。查询优化器倾向于采用聚集索引，因为聚集索引能够在 B+ 树索引的叶子节点上直接找到数据。

聚集索引对于主键的排序查找和范围查找速度非常快。

两点关键信息：

- 根据主键值创建了 B+ 树结构
- 每个叶子节点包含了整行数据

##### 辅助索引

上图中两点关键点需要注意：

- 根据 a 字段的值创建了 B+ 树结构
- 每个叶子节点保存的是 a 字段自己的键值和主键 ID

下面这条查询语句：

```sql
select * from t8 where a=3;
```

它先通过 a 字段上的索引树，得到主键 id 为 3，再到 id 的聚集索引树上找到对应的行数据。

而下面这条 SQL:

```sql
select * from t8 where id=3;
```

查询到的结果是一样的，而执行过程则只需要搜索 id 的聚集索引树。我们能看出辅助索引的查询比主键查询多扫描一颗索引树，所以，我们应该**尽量使用主键做为条件进行查询**。

### 2  哪些情况下需要添加索引？

目前比较常见需要创建索引的场景有：数据检索时在条件字段添加索引、聚合函数对聚合字段添加索引、对排序字段添加索引、为了防止回表添加索引、关联查询在关联字段添加索引等。

#### 数据检索

把有索引的字段 a 作为条件进行查询

```sql
select * from t9_1 where a = 90000;
```

执行计划中：

前者 type 字段为 ALL，后者 type 字段为 ref，显然后者性能更好

rows 这个字段前者是 1596288，而后者是 16，有索引的情况扫描行数大大降低。

因此建议数据检索时，在条件字段添加索引。

#### 聚合函数

求有索引的字段 a 的最大值：

```sql
select max(a) from t9_1;
```

相比对没有索引的字段 d 求最大值（花费330毫秒），**显然索引能提升 max() 函数的效率，同理也能提升 min() 函数的效率**。

从 MySQL 5.7.18 开始，通过遍历最小的可用二级索引来处理 count(*) 语句，如果不存在二级索引，则扫描聚簇索引。原因是：InnoDB 二级索引树的叶子节点上存放的是主键，而主键索引树的叶子节点上存放的是整行数据，所以二级索引树比主键索引树小。因此优化器基于成本的考虑，优先选择的是二级索引。

**因此索引对聚合函数 count(\*) 也有优化作用。**

#### 排序

通过添加合适索引优化 order by 的方法

- 如果对单个字段排序，则可以在这个排序字段上添加索引来优化排序语句；
- 如果是多个字段排序，可以在多个排序字段上添加联合索引来优化排序语句；
- 如果是先等值查询再排序，可以通过在条件字段和排序字段添加联合索引来优化排序语句。

####  避免回表

比如下面这条 SQL：

```sql
select a,d from t9_1 where a=90000;
```

可以走 a 字段的索引，但是在学了第 8 节后，我们知道了辅助索引的结构，如果通过辅助索引来寻找数据，InnoDB 存储引擎会遍历辅助索引树查找到对应记录的主键，然后通过主键索引回表去找对应的行数据。

但是，如果条件字段和需要查询的字段有联合索引的话，其实回表这一步就省了，因为联合索引中包含了这两个字段的值。像这种索引就已经覆盖了我们的查询需求的场景，我们称为：**覆盖索引**。比如下面这条 SQL：

```sql
select b,c from t9_1 where b=90000;
预览
```

可直接通过联合索引 idx_b_c 找到 b、c 的值

**所以可以通过添加覆盖索引让 SQL 不需要回表，从而减少树的搜索次数，让查询更快地返回结果。**

####  关联查询

### 3   普通索引和唯一索引有哪些区别？

对于普通索引和唯一索引的区别，也许你已经知道：有普通索引的字段可以写入重复的值，而有唯一索引的字段不可以写入重复的值。其实对于 MySQL 来说，不止这一种区别。

#### Insert Buffer

为什么要增加 Insert Buffer？

增加 Insert Buffer 有两个好处：

- 减少磁盘的离散读取
- 将多次插入合并为一次操作

但是得注意的是，使用 Insert Buffer 得满足两个条件：

- 索引是辅助索引
- 索引不是唯一

跟 Insert Buffer 一样，Change Buffer 也得满足这两个条件：

- 索引是辅助索引
- 索引不是唯一

#### Change Buffer

为什么唯一索引的更新不使用 Change Buffer ?

原因：唯一索引**必须要将数据页读入内存才能判断是否违反唯一性约束**。如果都已经读入到内存了，那直接更新内存会更快，就没必要使用 Change Buffer 了。

#### 普通索引和唯一索引的区别

- 数据修改时，普通索引可以用 Change Buffer，而唯一索引不行。
- 数据修改时，唯一索引在 RR 隔离级别下，更容易出现死锁。
- 查询数据是，普通索引查到满足条件的第一条记录还需要继续查找下一个记录，而唯一索引查找到第一个记录就可以直接返回结果了，但是普通索引多出的查找次数所消耗的资源多数情况可以忽略不计。

####  普通索引和唯一索引如何选择

如果业务要求某个字段唯一，但是代码不能完全保证写入唯一值，则添加唯一索引，让这个字段唯一，该字段新增重复数据时，将报类似如下的错：

```sql
ERROR 1062 (23000): Duplicate entry '1' for key 'f1'
```

如果代码确定某个字段不会有重复的数据写入，则可以选择添加普通索引。 因为普通索引可以使用 Change Buffer，并且出现死锁的概率比唯一索引低。

### 4  联合索引有哪些讲究？

#### 认识联合索引

联合索引：是指对表上的多个列进行索引。适合 where 条件中的多列组合，在某些场景可以避免回表。

联合索引的键值数量大于 1（比如上图中有 a 和 b 两个键值），与单个键值的 B+ 树一样，也是按照键值排序的。比如图中 a、b 两个字段的值为 (1,1),(1,2),(1,3),(2,1),(2,2),(2,3)，是按(a,b) 进行排序的。因此，**对于 a、b 两个字段都做为条件时，查询是可以走索引的；对于单独 a 字段查询也是可以走索引的。但是对于 b 字段单独查询就走不了索引了**，因为在上图，b 字段对应的值为 1,2,3,1,2,3，显然不是有序的，所以走不了 b 字段的索引。

联合索引的建议：

- where 条件中，经常同时出现的列放在联合索引中。
- 把选择性最大的列放在联合索引的最左边。

#### 联合索引使用分析

##### 可以完整用到联合索引的情况

下面我们列出几种可以完整用到联合索引的情况

```sql
select * from t11 where a=1 and b=1 and c=1; /* sql1 */
```

##### 只能使用部分联合索引的情况

有些场景只能用到部分联合索引，这里就列出几种情况。

```sql
select * from t11 where a=1 and b=1; /* sql11 */
```

对于联合索引 idx_a_b_c（a,b,c） ，如果条件中只包含 a 和 c，则只能用到联合索引中 a 的索引。c 这里是用不了索引的。**联合索引 idx_a_b_c(a,b,c) 相当于 (a) 、(a,b) 、(a,b,c) 三种索引，称为联合索引的最左原则。**

##### 可以用到覆盖索引的情况

什么是覆盖索引？

从辅助索引中就可以查询到结果，不需要回表查询聚集索引中的记录。

使用覆盖索引的优势：因为不需要扫描聚集索引，因此可以减少 SQL 执行过程的 IO 次数。

```sql
select b,c from t11 where a=3; /* sql21 */
```

通过 a 字段上的条件，去联合索引 idx_a_b_c 的索引树上可以直接查找到 b 字段和 c 字段的值，不需要回表，因此 sql21 使用到了覆盖索引。

通过 a 字段上的条件，去联合索引 idx_a_b_c 的索引树上可以直接查找到 b 字段和 c 字段的值，不需要回表，因此 sql21 使用到了覆盖索引。

```sql
select c from t11 where a=1 and b=1 ; /* sql22 */
```

跟 sql21 类似，在联合索引 idx_a_b_c 的索引树上，通过 a 和 b 的值可以直接找到 c 的值，因此 sql22 使用的也是覆盖索引。

```sql
select id from t11 where a=1 and b=1 and c=1; /* sql23 */
```

通过 a、b、c 三个字段的值，去联合索引树的叶子节点找到主键 id，不需要回表，因此 sql23 也使用了覆盖索引。

### 5  为什么MySQL会选错索引？

#### show index的使用

上面几个重要的字段做一下解释：

- Non_unique：如果是唯一索引，则值为 0，如果可以有重复值，则值为 1
- Key_name：索引名字
- Seq_in_index：索引中的列序号，比如联合索引 idx_a_b_c (a,b,c) ，那么三个字段分别对应 1,2,3
- Column_name：字段名
- Collation：字段在索引中的排序方式，A 表示升序，NULL 表示未排序
- Cardinality：索引中不重复记录数量的预估值，该值等会儿会详细讲解
- Sub_part：如果是前缀索引，则会显示索引字符的数量；如果是对整列进行索引，则该字段值为 NULL
- Null：如果列可能包含空值，则该字段为 YES；如果不包含空值，则该字段值为 ’ ’
- Index_type：索引类型，包括 BTREE、FULLTEXT、HASH、RTREE 等

#### Cardinality 取值

#### 统计信息不准确导致选错索引

在 MySQL 中，**优化器控制着索引的选择。一般情况下，优化器会考虑扫描行数、是否使用临时表、是否排序等因素，然后选择一个最优方案去执行 SQL 语句**。

#### 单次选取的数据量过大导致选错索引

## MySQL锁

### 1  全局锁和表锁什么场景会用到？

**如何保证数据并发访问的一致性、有效性呢？**

MySQL 中，锁就是协调多个用户或者客户端并发访问某一资源的机制，保证数据并发访问时的一致性和有效性。

根据加锁的范围，MySQL 中的锁可分为三类：

- 全局锁
- 表级锁
- 行锁

#### 全局锁

MySQL 全局锁会关闭所有打开的表，并使用全局读锁锁定所有表。其命令为：

```sql
FLUSH TABLES WITH READ LOCK;
```

简称：FTWRL，可以使用下面命令解锁：

```sql
UNLOCK TABLES;
```

进行 FTWRL 实验：

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| FLUSH TABLES WITH READ LOCK; Query OK, 0 rows affected (0.00 sec) |                                                              |
| select * from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** | select * from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** |
| insert into t14(a,b) values(2,2); ERROR 1223 (HY000): Can’t execute the query because you have a conflicting read lock **（报错）** | insert into t14(a,b) values(2,2);/* sql1 */ **（等待）**     |
| UNLOCK TABLES;                                               | insert into t14(a,b) values(2,2);/* sql1 */ Query OK, 1 row affected (5.73 sec) **（session1 解锁后，在等待的 sql1 马上执行成功）** |

上面的实验中，当 session1 执行 FTWRL 后，本线程 session1 和其它线程 session2 都可以查询，本线程和其它线程都不能更新。

原因是：**当执行 FTWRL 后，所有的表都变成只读状态，数据更新或者字段更新将会被阻塞。**

全局锁一般用在整个库（包含非事务引擎表）做备份（mysqldump 或者 xtrabackup）时。也就是说，在整个备份过程中，整个库都是只读的，其实这样风险挺大的。如果是在主库备份，会导致业务不能修改数据；而如果是在从库备份，就会导致主从延迟。

好在 mysqldump 包含一个参数 --single-transaction，可以在一个事务中创建一致性快照，然后进行所有表的备份。因此增加这个参数的情况下，备份期间可以进行数据修改。但是需要所有表都是事务引擎表。所以这也是建议使用 InnoDB 存储引擎的原因之一。

而对于 xtrabackup，可以分开备份 InnoDB 和 MyISAM，或者不执行 --master-data，可以避免使用全局锁。

#### 表级锁

表级锁有两种：表锁和元数据锁。

##### 表锁

表锁使用场景：

1. 事务需要更新某张大表的大部分或全部数据。如果使用默认的行锁，不仅事务执行效率低，而且可能造成其它事务长时间锁等待和锁冲突，这种情况下可以考虑使用表锁来提高事务执行速度；
2. 事务涉及多个表，比较复杂，可能会引起死锁，导致大量事务回滚，可以考虑表锁避免死锁。

其中表锁又分为表读锁和表写锁，命令分别是：

表读锁：

```sql
lock tables t14 read;
```

表写锁：

```sql
lock tables t14  write;
```

下面我们分别用实验验证表读锁和表写锁。

表读锁实验：

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| lock tables t14 read; Query OK, 0 rows affected (0.00 sec)   |                                                              |
| select id,a,b from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** | select id,a,b from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** |
| insert into t14(a,b) values(3,3); ERROR 1099 (HY000): Table ‘t14’ was locked with a READ lock and can’t be updated **（报错）** | insert into t14(a,b) values(3,3);/* sql2 */ **（等待）**     |
| unlock tables; Query OK, 0 rows affected (0.00 sec)          | insert into t14(a,b) values(3,3);/* sql2 */ Query OK, 1 row affected (10.97 sec) **（session1 解锁后，sql2 立马写入成功）** |

从上面的实验我们可以看出，在 session1 中对表 t14 加表读锁，session1 和 session2 都可以查询表 t14 的数据；而 session1 执行更新会报错，session2 执行更新会等待（直到 session1 解锁后才更新成功）。

总结：**对表执行 lock tables xxx read （表读锁）时，本线程和其它线程可以读，本线程写会报错，其它线程写会等待。**

我们再来看一下表写锁实验：

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| lock tables t14 write; Query OK, 0 rows affected (0.00 sec)  |                                                              |
| select id,a,b from t14 limit 1; … 1 row in set (0.00 sec) **（能正常返回结果）** | select id,a,b from t14 limit 1;/* sql3 */ **（等待）**       |
| unlock tables; Query OK, 0 rows affected (0.01 sec)          | select id,a,b from t14 limit 1;/* sql3 */ … 1 row in set (7.16 sec) **（session1 解锁后，sql3 马上返回查询结果）** |
| lock tables t14 write; Query OK, 0 rows affected (0.00 sec)  |                                                              |
| delete from t14 limit 1; Query OK, 1 row affected, 1 warning (0.00 sec) **（能正常执行删除语句）** | delete from t14 limit 1;/* sql4 */ **（等待）**              |
| unlock tables; Query OK, 0 rows affected (0.00 sec)          | delete from t14 limit 1;/* sql4 */ Query OK, 1 row affected, 1 warning (14.94 sec) **（session1 解锁后，sql4 立马执行成功）** |

总结：**对表执行 lock tables xxx write （表写锁）时，本线程可以读写，其它线程读写都会阻塞。**

##### 元数据锁

#### 总结

其中**全局锁会让所有的表变成只读状态，所有更新操作都会被阻塞。**

而表级锁分为表锁和元数据锁。

表锁又提到了表读锁和表写锁，并都进行了实验。两者的区别是：

**表读锁：本线程和其它线程可以读，本线程写会报错，其它线程写会等待。**

**表写锁：本线程可以读写，其它线程读写都会阻塞。**

为了保证事务和 DDl 并行执行数据一致，在 MySQL 5.5.3 引入了 MDL 锁。通过本节讲解的 MDL 锁机制，应该注意的几个点是：

- 尽量避免慢查询
- 事务要及时提交
- 避免大事务
- 避免在业务高峰执行 DDL 操作

### 2  行锁：InnoDB替代MyISAM的重要原因

MySQL 5.5 之前的默认存储引擎是 MyISAM，5.5 之后改成了 InnoDB。InnoDB 后来居上最主要的原因就是：

- InnoDB 支持事务：适合在并发条件下要求数据一致的场景。
- InnoDB 支持行锁：有效降低由于删除或者更新导致的锁定。

#### 两阶段锁

传统的关系型数据库加锁的一个原则是：两阶段锁原则。

两阶段锁：锁操作分为两个阶段，加锁阶段和解锁阶段，并且保证加锁阶段和解锁阶段不相交。

我们可以通过下面这张表理解两阶段锁：

| 序号 | MySQL 操作      | 解释                                    | 锁阶段   |
| :--- | :-------------- | :-------------------------------------- | :------- |
| 1    | begin;          | 事务开始                                |          |
| 2    | insert into …;  | 加 insert 对应的锁                      | 加锁阶段 |
| 3    | update table …; | 加 update 对应的锁                      | 加锁阶段 |
| 4    | delete from …;  | 加 delete 对应的锁                      | 加锁阶段 |
| 5    | commit;         | 事务结束，同时释放 2、3、4 步骤中加的锁 | 解锁阶段 |

####  InnoDB 行锁模式

InnoDB 实现了以下两种类型的行锁：

- 共享锁（S）：允许一个事务去读一行，阻止其它事务获得相同数据集的排他锁；
- 排他锁（X）：允许获得排他锁的事务更新数据，阻止其它事务取得相同数据集的共享读锁和排他写锁。

对于普通 select 语句，InnoDB 不会加任何锁，事务可以通过以下语句显式给记录集加共享锁或排他锁：

- 共享锁（S）：select * from table_name where … lock in share mode;
- 排他锁（X）：select * from table_name where … for update。

#### InnoDB 行锁算法

#### 事务隔离级别

MySQL 的 4 种隔离级别：

- Read uncommitted（读未提交）: 在该隔离级别，所有事务都可以看到其它未提交事务的执行结果。可能会出现脏读。
- Read Committed（读已提交，简称： RC）：一个事务只能看见已经提交事务所做的改变。因为同一事务的其它实例在该实例处理期间可能会有新的 commit，所以可能出现幻读。
- Repeatable Read（可重复读，简称：RR）：这是 MySQL 的默认事务隔离级别，它确保同一事务的多个实例在并发读取数据时，会看到同样的数据行。消除了脏读、不可重复读，默认也不会出现幻读。
- Serializable（串行）：这是最高的隔离级别，它通过强制事务排序，使之不可能相互冲突，从而解决幻读问题。

> 这里解释一下脏读和幻读：
>
> - 脏读：读取未提交的事务。
> - 幻读：一个事务按相同的查询条件重新读取以前检索过的数据，却发现其他事务插入了满足其查询条件的新数据。

#### RC 隔离级别下的行锁实验

##### 通过非索引字段查询

要想分析某条 SQL 是怎么加锁的，如果其他信息都不知道，那就得分几种情况了，不同情况加锁的方式也各不一样，比较常见的一些情况如下：

- RC 隔离级别，a 字段没索引。
- RC 隔离级别，a 字段有唯一索引。
- RC 隔离级别，a 字段有非唯一索引。
- RR 隔离级别，a 字段没索引。
- RR 隔离级别，a 字段有唯一索引。
- RR 隔离级别，a 字段有非唯一索引。

**没有索引的情况下，InnoDB 的当前读会对所有记录都加锁。所以在工作中应该特别注意 InnoDB 这一特性，否则可能会产生大量的锁冲突。**

##### 通过唯一索引查询

**如果查询的条件是唯一索引，那么 SQL 需要在满足条件的唯一索引上加锁，并且会在对应的聚簇索引上加锁。**

##### 通过非唯一索引查询

**如果查询的条件是非唯一索引，那么 SQL 需要在满足条件的非唯一索引上都加上锁，并且会在它们对应的聚簇索引上加锁。**

### 3  答疑篇

#### 关于条件字段做函数操作的问题

下面这条 SQL 走不了索引：

```mysql
select * from t1 where data(c) = '2019-05-21';
```

原因是，索引树中存储的是列的实际值和主键值。如果拿 ‘2019-05-21’ 去匹配，将无法定位到索引树中的值。

有同学就问，如果把语句 “=” 前面的 DATE_FORMAT 按照实际值转换，比如：

```sql
select * from t1 where DATE_FORMAT(c,'%Y-%m-%d %H:%i:%s')='2019-05-21 00:00:00';
```

已经转换成跟实际值一样，为什么也不走索引？

原因是：这条 SQL 执行过程是先找出所有 c 字段的值，然后经过 DATE_FORMAT 加工后形成新的值集合（并且是无序的），因此最终是这个新值集合跟 ‘2019-05-21’ 进行比较。而索引中存储的是 c 字段的集合，所以无法使用索引过滤。

#### 是不是所有的辅助索引都需要回表？

并不是所有的辅助索引查询都需要回表，如果条件字段跟需要查询的字段在一颗索引树上，可以直接在索引树上找到结果，就不用回表了，通常称为覆盖索引。

#### Btree 索引和 HASH 索引有什么区别？

两者区别：

- Btree 索引可能需要多次运用折半查找来找到对应的数据库；
- 而 HASH 索引是通过 HASH 函数，计算出 HASH 值，在表中找出对应的数据。

优缺点对比：

- 大量不同数据等值精确查询，HASH 索引效率通常比 B+TREE 高；
- 但是 HASH 索引不支持模糊查询、排序、范围查询和联合索引中的最左匹配原则，而这些 Btree 索引都支持。

#### 关于隐式转换的一个问题

表 t1 中 a 字段是字符串类型，有索引，但是如下 SQL 走不了索引：

```sql
select * from t1 where a=1000;
```

原因是：MySQL 会把 a 转换成 int 型，再去做判断，也相当于对索引字段做函数操作，导致不能使用索引。

有同学提问：表 t1 中，为 int 型的 b 字段，有索引，下面 SQL 为什么可以走索引？

```sql
select * from t1 where b='1';
```

因为：MySQL 中如果条件字段和后面的值不是同一类型，那么优先将字符串转成数字类型做比较。因此上面的 SQL 是转的 “=” 后面的值，而不是字段 b，因此 SQL 其实等价于：

```sql
select * from t1 where b=1;
```

### 4  间隙锁的意义

### 5  为什么会出现死锁？

**死锁是指两个或者多个事务在同一资源上相互占用，并请求锁定对方占用的资源，从而导致恶性循环的现象。**

InnoDB 中解决死锁问题有两种方式：

1. 检测到死锁的循环依赖，立即返回一个错误（这个报错内容请看下面的实验），将参数 innodb_deadlock_detect 设置为 on 表示开启这个逻辑；
2. 等查询的时间达到锁等待超时的设定后放弃锁请求。这个超时时间由 innodb_lock_wait_timeout 来控制。默认是 50 秒。

#### 同一张表中

不同线程并发访问同一张表的多行数据，未按顺序访问导致死锁。

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| begin;                                                       | begin;                                                       |
| select * from t18 where a=1 for update; … 1 row in set (0.00 sec) | select * from t18 where a=2 for update; … 1 row in set (0.00 sec) |
| select * from t18 where a=2 for update;/* SQL1 */ （等待）   |                                                              |
| **（session2 提示死锁回滚后，SQL1 成功返回结构）**           | select * from t18 where a=1 for update; **ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting transaction** |
| commit;                                                      | commit;                                                      |

####  不同表之间

不同线程并发访问多个表时，未按顺序访问导致死锁：

| session1                                                     | session2                                                     |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| begin;                                                       | begin;                                                       |
| select * from t18 where a=1 for update; … 1 row in set (0.00 sec) | select * from t18_1 where a=1 for update; … 1 row in set (0.00 sec) |
| select * from t18_1 where a=1 for update;/* SQL2 */ 等待     |                                                              |
| **（session2 提示死锁回滚后，SQL1 成功返回结构）**           | select * from t18 where a=1 for update; **ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting transaction** |
| commit;                                                      | commit;                                                      |

#### 如何降低死锁概率？

那么应该怎样降低出现死锁的概率呢？这里总结了如下一些经验：

1. 更新 SQL 的 where 条件尽量用索引；
2. 基于 primary 或 unique key 更新数据；
3. 减少范围更新，尤其非主键、非唯一索引上的范围更新；
4. 加锁顺序一致，尽可能一次性锁定所有需要行；
5. 将 RR 隔离级别调整为 RC 隔离级别。

#### 分析死锁的方法

nnoDB 中，可以使用 SHOW INNODB STATUS 命令来查看最后一个死锁的信息。我们可以尝试用下这个命令获取一些死锁信息，如下：

```sql
show engine innodb status\G
```

## 事务

### 1  数据库忽然断电会丢数据吗？

#### 什么是事务？

事务就是一组原子性的 SQL 查询，
