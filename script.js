import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 10 },
    { duration: '2m', target: 200 },
    { duration: '1m', target: 10 },
  ],
};

export default function () {
  http.get('http://localhost:8090/health');
  sleep(1);
}
