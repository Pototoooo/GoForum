import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 10 },   // 逐步增加到 10 个并发
    { duration: '30s', target: 50 },   // 增加到 50 个并发
    { duration: '10s', target: 0 },    // 逐步减少到 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% 请求应在 500ms 内
    http_req_failed: ['rate<0.01'],    // 失败率低于 1%
  },
};

// 先通过列表接口获取一批帖子 ID
function fetchPostIds() {
  const res = http.get('http://localhost:8084/posts2?page=1&size=20&order=time');
  check(res, { 'list status 200': (r) => r.status === 200 });
  if (res.status !== 200) return [];
  const body = JSON.parse(res.body);
  if (body.code !== 1000 || !body.data) return [];
  return body.data.map(p => p.id);
}

export default function () {
  // 测试帖子列表接口
  const listRes = http.get('http://localhost:8084/posts2?page=1&size=10&order=time');
  check(listRes, {
    '列表接口状态 200': (r) => r.status === 200,
    '列表接口返回成功': (r) => JSON.parse(r.body).code === 1000,
  });

  sleep(1);
}