<?php
    include "xuncache.class.php";

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

 ?>
