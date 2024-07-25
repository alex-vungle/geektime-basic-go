import redis
import random

# 连接到本地的Redis服务器
client = redis.Redis(host='localhost', port=6379, db=0)

# ZSET的键
zset_key = 'top_100_article'

# 插入100个数据
for i in range(100):
    article_id = 100 + i
    score = random.randint(200, 2000)
    client.zadd(zset_key, {article_id: score})

print("Inserted 100 articles into the ZSET with random scores.")