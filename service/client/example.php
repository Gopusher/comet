<?php

class JsonRPC
{
    /**
     * redis实例
     */
    private static $obj = [];

    private function getConnect($host, $port)
    {
        $key = $host . ':' . $port;
        if (!isset(self::$obj[$key])) {
            $conn = fsockopen($host, $port, $errno, $errstr, 3);
            if (! $conn) {
                throw new \Exception('建立连接失败');
            }
            self::$obj[$key] = $conn;
        }
        return self::$obj[$key];
    }

    public function execute($host, $port, $data) {
        $conn = $this->getConnect($host, $port);

        if (! $conn) {
            throw new \Exception('连接不存在');
        }
        $err = fwrite($conn, json_encode(array_merge($data, ['id' => 0])) . "\n");

        if ($err === false) {
            throw new \Exception('调用方法失败' . __METHOD__);
        }

        stream_set_timeout($conn, 0, 3000);
        $line = fgets($conn);
        if ($line === false) {
            return NULL;
        }
        return json_decode($line,true);
    }

    public function SendToConnections($host, $port, $to, $msg)
    {
        $data = array(
            'method' => "Server.SendToConnections",
            'params' => [json_encode([
                'to'    => (array) $to,
                'msg'   => $msg
            ])],
        );

        return $this->execute($host, $port, $data);
    }

    public function __construct()
    {
        foreach (self::$obj as $conn) {
            fclose($conn);
        }
    }
}

if (count($argv) < 3) {
    throw new Exception('消息发送参数格式: php ./example.php msg ...to');
}

$s = time();
//todo 这里的 ip 要注意
$client = new JsonRPC();
//'{ "from": "5ab22da3d237a", "to": "1", "type": "group", "contentType": "text", "content": "123" }'
for ($i = 0; $i < 20000; $i++) {
    $r = $client->sendToConnections("10.0.1.131", 8901, array_slice($argv, 2), $argv[1]);
    // var_export($r);
}
$e = time();
echo $e - $s;
