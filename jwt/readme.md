# jwt

tkjwt 需要redis作为缓存，所以：请初始化redis连接后再调用

```text

// redis 初始化
tkredis.Setup("redis.toml", "Client")
tkredis.Close()

```