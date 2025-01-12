# learning-redis

- [x] Create, Read, Update, Delete calls in Redis (SET, GET, DEL)
- [x] Basic Pub/Sub implementation (Publish, Subscribe)
- [x] Event Logs (Rpush, LRange)
- [x] LeaderBoard (ZAdd, ZRevRangeWithScores)
- [x] Counter (Inc / Using some Gin too)
- [x] Likes of a given post (SAdd, SMembers)
- [x] Redis Hash (HSET, HGETALL, HLEN)
- [x] Real Time notifications (using Server Sent Events, Publish, Subscribe)
- [x] Temporary URL (using SET, and TTL)
- [x] Uploading image and storing as bas64 (SET, GET)

Running Redis using Docker

```
docker compose up -d
```

Shutting down Redis using Docker

```
docker compose down
```