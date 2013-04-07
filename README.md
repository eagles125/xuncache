## xuncache
========
xuncache 是免费开源的NOSQL(内存数据库) 采用golang开发,简单易用而且 功能强大(就算新手也完全胜任)、性能卓越能轻松处理海量数据,可用于缓存系统.

目前版本 version 0.2

采用json协议 socket通信 --后期打算用bson

## 目前功能
========
-增加or设置

-查找数据

-删除数据

-暂不支持key过期操作

支持 php 客户端

LICENSE: under the BSD License
- by [孙彦欣](http://weibo.com/sun8911879)
## php代码示例
========

	$xuncache = new xuncache();
    $data['name'] = "xuncache";
    $data['time'] = "20130408";
    //添加数据
    $status = $xuncache->key("syx")->add($data);
    dump($status);
    //查找数据
    $cache = $xuncache->key("syx")->find();
    dump($cache);
    //删除数据
    $status = $xuncache->key("syx")->del();
    dump($status);
	//////////返回数据
	bool(true)
	array(2) {
	  ["name"] => string(8) "xuncache"
	  ["time"] => string(8) "20130408"
	}
	bool(true)
