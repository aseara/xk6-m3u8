import m3u8 from 'k6/x/m3u8';

export const options = {
    vus: 2,
    duration: '30s',
};

m3u8.set("http://127.0.0.1:30769/01.m3u8")

export default function () {
    m3u8.record()
}
