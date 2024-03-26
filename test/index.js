import http from 'k6/http';
import { check, group, sleep } from 'k6';

export const options = {
  vus: 1,
}

export default function() {
  const SERVER_URL = "http://localhost:8080"
  const Headers = {
    'Content-Type': 'application/json',
  }

  let res = http.get(`${SERVER_URL}/ping`);
  check(res, { 'Server can ping': (r) => r.status == 200 });

  res = http.get(`${SERVER_URL}/healthy`);
  check(res, { 'Server is healthy': (r) => r.status == 200 });
}
