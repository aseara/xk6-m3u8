# xk6-redis

This is a m3u8 download library for [k6](https://github.com/grafana/k6),
implemented as an extension using the [xk6](https://github.com/grafana/xk6) system.

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Install `xk6`:
  ```shell
  go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```shell
  xk6 build --with github.com/grafana/xk6-m3u8
  ```

## Usage

This extension exposes a [promise](https://javascript.info/promise-basics)-based API. As opposed to most other current k6 modules and extensions, who operate in a synchronous manner,
xk6-m3u8 operates in an asynchronous manner. In practice, this means that using the Redis client's methods won't block the execution of the test,
and that the test will continue to run even if the Redis client is not ready to respond to the request.

One potential caveat introduced by this new execution model, is that to depend on operations against Redis, all following operations
should be made in the context of a [promise chain](https://javascript.info/promise-chaining). As demonstrated in the examples below, whenever there is a dependency on a Redis
operation result or return value, the operation, should be wrapped in a promise itself. That way, a user can perform asynchronous interactions
with Redis in a seemingly synchronous manner.

For instance, if you were to depend on values stored in Redis to perform HTTP calls, those HTTP calls should be made in the context of the Redis promise chain:

```javascript
// Instantiate a new redis client
const redisClient = new redis.Client({
  addrs: __ENV.REDIS_ADDRS.split(",") || new Array("localhost:6379"), // in the form of "host:port", separated by commas
  password: __ENV.REDIS_PASSWORD || "",
})

export default function() {
    // Once the SRANDMEMBER operation is succesfull,
    // it resolves the promise and returns the random
    // set member value to the caller of the resolve callback.
    //
    // The next promise performs the synchronous HTTP call, and
    // returns a promise to the next operation, which uses the 
    // passed URL value to store some data in redis.
    redisClient.srandmember('client_ids')
        .then((randomID) => {
            const url = `https://my.url/${randomID}`
            const res = http.get(url)

            // Do something with the result...

            // return a promise resolving to the URL
            return url
        })
        .then((url) => redisClient.hincrby('k6_crocodile_fetched', url, 1))
}
```

## Example test scripts

In this example we demonstrate two scenarios: one load testing a redis instance, another using redis as an external data store used throughout the test itself.

```javascript
import { check } from "k6";
import http from "k6/http";
import redis from "k6/x/redis";
import exec from "k6/execution";
import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.1/index.js";

export const options = {
  scenarios: {
    redisPerformance: {
      executor: "shared-iterations",
      vus: 10,
      iterations: 200,
      exec: "measureRedisPerformance",
    },
    usingRedisData: {
      executor: "shared-iterations",
      vus: 10,
      iterations: 200,
      exec: "measureUsingRedisData",
    },
  },
};

// Instantiate a new redis client
const redisClient = new redis.Client({
  addrs: __ENV.REDIS_ADDRS.split(",") || new Array("localhost:6379"), // in the form of "host:port", separated by commas
  password: __ENV.REDIS_PASSWORD || "",
});

// Prepare an array of crocodile ids for later use
// in the context of the measureUsingRedisData function.
const crocodileIDs = new Array(0, 1, 2, 3, 4, 5, 6, 7, 8, 9);

export function measureRedisPerformance() {
  // VUs are executed in a parallel fashion,
  // thus, to ensure that parallel VUs are not
  // modifying the same key at the same time,
  // we use keys indexed by the VU id.
  const key = `foo-${exec.vu.idInTest}`;

  redisClient
    .set(`foo-${exec.vu.idInTest}`, 1)
    .then(() => redisClient.get(`foo-${exec.vu.idInTest}`))
    .then((value) => redisClient.incrBy(`foo-${exec.vu.idInTest}`, value))
    .then((_) => redisClient.del(`foo-${exec.vu.idInTest}`))
    .then((_) => redisClient.exists(`foo-${exec.vu.idInTest}`))
    .then((exists) => {
      if (exists !== 0) {
        throw new Error("foo should have been deleted");
      }
    });
}

export function setup() {
  redisClient.sadd("crocodile_ids", ...crocodileIDs);
}

export function measureUsingRedisData() {
  // Pick a random crocodile id from the dedicated redis set,
  // we have filled in setup().
  redisClient
    .srandmember("crocodile_ids")
    .then((randomID) => {
      const url = `https://test-api.k6.io/public/crocodiles/${randomID}`;
      const res = http.get(url);

      check(res, {
        "status is 200": (r) => r.status === 200,
        "content-type is application/json": (r) =>
          r.headers["content-type"] === "application/json",
      });

      return url;
    })
    .then((url) => redisClient.hincrby("k6_crocodile_fetched", url, 1));
}

export function teardown() {
  redisClient.del("crocodile_ids");
}

export function handleSummary(data) {
  redisClient
    .hgetall("k6_crocodile_fetched")
    .then((fetched) => Object.assign(data, { k6_crocodile_fetched: fetched }))
    .then((data) =>
      redisClient.set(`k6_report_${Date.now()}`, JSON.stringify(data))
    )
    .then(() => redisClient.del("k6_crocodile_fetched"));

  return {
    stdout: textSummary(data, { indent: "  ", enableColors: true }),
  };
}
```

Result output:

```shell
$ ./k6 run test.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: test.js
     output: -

  scenarios: (100.00%) 1 scenario, 10 max VUs, 1m30s max duration (incl. graceful stop):
           * default: 10 looping VUs for 1m0s (gracefulStop: 30s)


running (1m00.1s), 00/10 VUs, 4954 complete and 0 interrupted iterations
default ✓ [======================================] 10 VUs  1m0s
```
