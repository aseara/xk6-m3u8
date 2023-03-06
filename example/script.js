import m3u8 from 'k6/x/m3u8';

export const options = {
    vus: 3,
    duration: '60s',
};



export function setup() {
    m3u8.start("http://127.0.0.1:30769/01.m3u8", ".\\tmp\\")
}

export default function () {
    m3u8.check()
}

export function teardown() {
    m3u8.stop()
}