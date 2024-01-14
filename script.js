import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  // Key configurations for spike in this section
  stages: [
    // { duration: '1m', target: 0 },
    { duration: '2m', target: 2000 }, // fast ramp-up to a high point
    { duration: '1m', target: 0 },
  ],
};

export default function () {
  http.get('http://localhost:3000/health');
  sleep(1);
}
