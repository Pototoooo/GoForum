import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// 自定义指标
const readFailureRate = new Rate('read_failed');
const writeFailureRate = new Rate('write_failed');
const postListTrend = new Trend('post_list_duration');
const postDetailTrend = new Trend('post_detail_duration');
const createPostTrend = new Trend('create_post_duration');
const loginTrend = new Trend('login_duration');

const BASE_URL = 'http://localhost:8084/api/v1';

// 预注册的测试用户
const TEST_USER = { username: 'k6user', password: 'password123' };

export const options = {
  stages: [
    { duration: '10s', target: 10 },   // 预热到 10 并发
    { duration: '30s', target: 50 },   // 升到 50 并发
    { duration: '20s', target: 100 },  // 升到 100 并发
    { duration: '10s', target: 0 },    // 降回 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'avg<200'],
    http_req_failed: ['rate<0.01'],
    read_failed: ['rate<0.01'],
    write_failed: ['rate<0.01'],
  },
};

// ====== 第一阶段：读接口基准测试 ======
export function readBenchmark() {
  // 帖子列表
  {
    const res = http.get(`${BASE_URL}/posts2?page=1&size=10&order=time`);
    postListTrend.add(res.timings.duration);
    const ok = check(res, {
      '列表状态200': r => r.status === 200,
      '列表返回成功': r => JSON.parse(r.body).code === 1000,
    });
    if (!ok) readFailureRate.add(1);
  }

  // 帖子详情 — 随机取一个存在的帖子
  {
    const res = http.get(`${BASE_URL}/post/453377581601787904`);
    postDetailTrend.add(res.timings.duration);
    const ok = check(res, {
      '详情状态200': r => r.status === 200,
    });
    if (!ok) readFailureRate.add(1);
  }

  // 社区列表
  {
    const res = http.get(`${BASE_URL}/community`);
    check(res, { '社区状态200': r => r.status === 200 });
  }
}

// ====== 第二阶段：写接口测试（需要登录） ======
export function writeBenchmark(token) {
  // 发帖
  {
    const payload = JSON.stringify({
      title: `负载测试帖子_${__VU}_${__ITER}`,
      content: 'k6 负载测试内容',
      community_id: 1,
    });
    const res = http.post(`${BASE_URL}/post`, payload, {
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    });
    createPostTrend.add(res.timings.duration);
    const ok = check(res, {
      '发帖状态200': r => r.status === 200,
    });
    if (!ok) writeFailureRate.add(1);
  }
}

// ====== 第三阶段：登录（获取 token） ======
export function login() {
  const res = http.post(`${BASE_URL}/login`,
    JSON.stringify(TEST_USER),
    { headers: { 'Content-Type': 'application/json' } }
  );
  loginTrend.add(res.timings.duration);
  check(res, { '登录状态200': r => r.status === 200 });
  if (res.status === 200) {
    try {
      return JSON.parse(res.body).data.token;
    } catch (e) {
      return '';
    }
  }
  return '';
}

// ====== setup：先注册测试用户 ======
export function setup() {
  // 注册用户（如果已存在会返回错误，没关系）
  http.post(`${BASE_URL}/signup`,
    JSON.stringify({ ...TEST_USER, re_password: TEST_USER.password }),
    { headers: { 'Content-Type': 'application/json' } }
  );

  // 登录获取 token
  const loginRes = http.post(`${BASE_URL}/login`,
    JSON.stringify(TEST_USER),
    { headers: { 'Content-Type': 'application/json' } }
  );

  let token = '';
  if (loginRes.status === 200) {
    try { token = JSON.parse(loginRes.body).data.token; } catch (e) {}
  }
  return { token };
}

// ====== default：混合场景 ======
export default function (data) {
  const r = Math.random();

  // 80% 读操作
  if (r < 0.8) {
    readBenchmark();
  }
  // 15% 写操作
  else if (r < 0.95) {
    writeBenchmark(data.token);
  }
  // 5% 登录
  else {
    login();
  }

  sleep(1);
}
