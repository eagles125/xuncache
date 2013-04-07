<?php
	/**
	 +------------------------------------------------------------------------------
	 * xuncache模型类
	 +------------------------------------------------------------------------------
	 * @author    sun8911879 <joijoi360@gmail.com>
	 * @version   $Id: xuncache.class.php 2013-04-06 22:54:00 sun8911879 $
	 +------------------------------------------------------------------------------
	 */
	class xuncache {
		// xuncache操作对象
    	protected $xuncache;
    	//静态实例化对象
    	protected static $socket;
        //连接IP
        protected static $addr = "127.0.0.1";
        //连接端口
        protected static $port = "3351";
        //连接密码
        protected static $password = "";
        // 调试模式
        protected static $debug = true;
        // 查询表达式参数
        protected static $options = array();

	/**
     +----------------------------------------------------------
     * 架构函数
     * 取得DB类的实例对象
     +----------------------------------------------------------
     * @param string $name 模型名称
     * @param mixed $connection 数据库连接信息
     +----------------------------------------------------------
     * @access public
     +----------------------------------------------------------
     */
    public function __construct() {
        
        // 数据库初始化操作
        // 当前模型数据库连接信息
        if(!self::$socket){
            self::$socket = @socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
            if (self::$socket < 1) {
                return false;
            }
            $source = @socket_connect(self::$socket, self::$addr, self::$port);
            if ($source < 1) {
                if(self::$debug == true){
                    $this->customError("connect","Unable to connect the server side");
                }else{
                    return false;
                }
            }
        }

    }

    /**
     *----------------------------------------------------------
     * 利用__call方法实现一些特殊的Model方法
     *----------------------------------------------------------
     * @access public
     *----------------------------------------------------------
     * @param string $method 方法名称
     * @param array $args 调用参数
     *----------------------------------------------------------
     * @return mixed
     *----------------------------------------------------------
     */
    public function __call($method,$args) {
        $method = strtolower($method);
        // 连贯操作的实现
        if(in_array($method,array('key'),true)) {
            self::$options[$method] = $args[0];
            return $this;
        }else{
            return false;
        }
    }

    /**
     *----------------------------------------------------------
     * 查询数据
     *----------------------------------------------------------
     * @access public
     *----------------------------------------------------------
     * @param mixed $options 表达式参数
     *----------------------------------------------------------
     * @return mixed
     *----------------------------------------------------------
     */
    public function find(){
        $array['Pass'] = self::$password;
        $array['Key'] = self::$options['key'];
        $array['Protocol'] = "find";
        $arrange = json_encode($array);
        if(!@socket_write(self::$socket, $arrange, strlen($arrange))){
            return false;
        }else{

            $accept = @socket_read(self::$socket, 8192);
            $accept = json_decode($accept);
            //状态判断
            if(@$accept->error == true&&self::$debug == true){
                $this->customError("connect",@$accept->point);
            }elseif(@$accept->error == true){
                return false;
            }
            if(count((array)$accept->data)<1){
                return false;
            }
        }
        return (array)$accept->data;
    }

    /**
     *----------------------------------------------------------
     * 添加数据
     *----------------------------------------------------------
     * @access public
     *----------------------------------------------------------
     * @param mixed $options 表达式参数
     *----------------------------------------------------------
     * @return mixed
     *----------------------------------------------------------
     */
    public function add($arr){
        $array['Pass'] = self::$password;
        $array['Key'] = self::$options['key'];
        $array['Protocol'] = "set";
        $array['Data'] = $arr;
        $arrange = json_encode($array);
        if(!@socket_write(self::$socket, $arrange, strlen($arrange))){
            return false;
        }else{
            $accept = @socket_read(self::$socket, 8192);
            $accept = json_decode($accept);
            //状态判断
            if(@$accept->error == true&&self::$debug == true){
                $this->customError("connect",@$accept->point);
            }elseif(@$accept->error == true){
                return false;
            }
            
        }
        return (bool)$accept->status;
    }

    /**
     *----------------------------------------------------------
     * 删除数据
     *----------------------------------------------------------
     * @access public
     *----------------------------------------------------------
     * @param mixed $options 表达式参数
     *----------------------------------------------------------
     * @return mixed
     *----------------------------------------------------------
     */
    public function del(){
        $array['Pass'] = self::$password;
        $array['Key'] = self::$options['key'];
        $array['Protocol'] = "delete";
        $arrange = json_encode($array);
        if(!@socket_write(self::$socket, $arrange, strlen($arrange))){
            return false;
        }else{
            $accept = @socket_read(self::$socket, 8192);
            $accept = json_decode($accept);
            //状态判断
            if(@$accept->error == true&&self::$debug == true){
                $this->customError("connect",@$accept->point);
            }elseif(@$accept->error == true){
                return false;
            }
            
        }
        return (bool)$accept->status;
    }



    private function customError($errno, $errstr){ 
        echo "<b>Error:</b> [$errno] $errstr<br />";
        die();
    }
}  


    // 浏览器友好的变量输出
    function dump($var, $echo=true, $label=null, $strict=true) {
        $label = ($label === null) ? '' : rtrim($label) . ' ';
        if (!$strict) {
            if (ini_get('html_errors')) {
                $output = print_r($var, true);
                $output = "<pre>" . $label . htmlspecialchars($output, ENT_QUOTES) . "</pre>";
            } else {
                $output = $label . print_r($var, true);
            }
        } else {
            ob_start();
            var_dump($var);
            $output = ob_get_clean();
            if (!extension_loaded('xdebug')) {
                $output = preg_replace("/\]\=\>\n(\s+)/m", "] => ", $output);
                $output = '<pre>' . $label . htmlspecialchars($output, ENT_QUOTES) . '</pre>';
            }
        }
        if ($echo) {
            echo($output);
            return null;
        }else
            return $output;
    }
?>