import redis from 'k6/x/redis';

const client = redis.newClient("192.168.2.39:26379", "Hg#9tk^s", 6);

export function setup() {
    redis.set(client,"snake","camel",0)
    redis.set(client,"foo",100,10)
}

export default function () {
    console.log(redis.get(client,"snake"))
    console.log(redis.get(client,"foo"))
    if (redis.do(client,"PING","bzzz") == "bzzz"){
        console.log("PONG!")
    }
}

export function teardown () {
    redis.del(client,"foo")
}