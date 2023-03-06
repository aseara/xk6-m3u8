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
  xk6 build --with github.com/grafana/xk6-m3u8=.
  ```

## Example test scripts


```javascript
import m3u8 from 'k6/x/m3u8';

const player = m3u8.Client();

export function setup() {
    player.start("http://127.0.0.1:30769/01.m3u8")
}

export default function () {
    player.check()
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
