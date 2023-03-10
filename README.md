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
  xk6 build --with github.com/aseara/xk6-m3u8=.
  ```

## Example test scripts


```javascript
import m3u8 from 'k6/x/m3u8';

export const options = {
    vus: 1,
    duration: '20s',
};

m3u8.set("http://127.0.0.1:30769/01.m3u8")

export default function () {
    m3u8.record()
}
```

Result output:

```shell
$ ./k6 run example/script.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: .\example\script.js
     output: engine

  scenarios: (100.00%) 1 scenario, 1 max VUs, 50s max duration (incl. graceful stop):
           * default: 1 looping VUs for 20s (gracefulStop: 30s)

INFO[0000] 2023/03/10 14:03:03 Start record live streaming movie at tmp\1\video.ts...
INFO[0000] 2023/03/10 14:03:03 Recorded segment  2023031003138                                                                                                       
INFO[0000] 2023/03/10 14:03:03 Recorded segment  2023031003139
INFO[0000] 2023/03/10 14:03:03 Recorded segment  2023031003140
INFO[0010] 2023/03/10 14:03:13 Start record live streaming movie at tmp\1\video.ts...
INFO[0010] 2023/03/10 14:03:13 Recorded segment  2023031003141
INFO[0010] 2023/03/10 14:03:13 Recorded segment  2023031003142

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=10.51s min=10.01s med=10.51s max=11s p(90)=10.9s p(95)=10.95s
     iterations...........: 2   0.095146/s
     vus..................: 1   min=1      max=1
     vus_max..............: 1   min=1      max=1

                                                                                                                                                                     
running (21.0s), 0/1 VUs, 2 complete and 0 interrupted iterations                                                                                                    
default ✓ [======================================] 1 VUs  20s 
```
