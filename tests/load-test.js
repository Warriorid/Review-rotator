import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 4 },
    { duration: '1m', target: 4 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'],
    http_req_failed: ['rate<0.001'],
    http_reqs: ['rate>5'],
  },
};

const baseURL = 'http://localhost:8080';

let testData = {
  teams: [],
  users: [],
  prs: []
};

export function setup() {
  const data = createTestData();
  return data;
}

export default function(data) {
  testData = data;
  
  if (testData.teams.length === 0) {
    return;
  }
  
  performOperations();
}

function createTestData() {
  const teams = [];
  const users = [];
  const prs = [];
  const MAX_TEAMS = 8;

  for (let teamNumber = 0; teamNumber < MAX_TEAMS; teamNumber++) {
    const teamName = `team-load-${teamNumber}`;
    
    const teamPayload = {
      team_name: teamName,
      members: [
        { user_id: `user-load-${teamNumber}-0`, username: `User ${teamNumber}-0`, is_active: true },
        { user_id: `user-load-${teamNumber}-1`, username: `User ${teamNumber}-1`, is_active: true },
        { user_id: `user-load-${teamNumber}-2`, username: `User ${teamNumber}-2`, is_active: true },
        { user_id: `user-load-${teamNumber}-3`, username: `User ${teamNumber}-3`, is_active: true }
      ]
    };
    
    const teamResponse = http.post(`${baseURL}/team/add`, JSON.stringify(teamPayload), {
      headers: { 'Content-Type': 'application/json' },
      timeout: '5s'
    });
    
    if (teamResponse.status === 201) {
      teams.push(teamName);
      teamPayload.members.forEach(member => users.push(member));
      
      for (let prNum = 0; prNum < 3; prNum++) {
        const prId = `pr-load-${teamNumber}-${prNum}`;
        const prPayload = {
          pull_request_id: prId,
          pull_request_name: `Load Test PR ${prId}`,
          author_id: `user-load-${teamNumber}-${prNum}`
        };
        
        const prResponse = http.post(`${baseURL}/pullRequest/create`, JSON.stringify(prPayload), {
          headers: { 'Content-Type': 'application/json' },
          timeout: '5s'
        });
        
        if (prResponse.status === 201) {
          prs.push(prId);
        }
      }
    }
    
    sleep(0.08);
  }
  
  return { teams, users, prs };
}

function performOperations() {
  const operations = [
    getTeam,
    getStats, 
    getUserReviews,
    getTeamAgain,
    createPR,
    getStatsAgain
  ];
  
  const numOperations = 3 + Math.floor(Math.random() * 2);
  
  for (let i = 0; i < numOperations; i++) {
    const randomOp = operations[Math.floor(Math.random() * operations.length)];
    randomOp();
    
    if (i < numOperations - 1) {
      sleep(0.015);
    }
  }
  
  sleep(0.08 + Math.random() * 0.08);
}

function getTeam() {
  const teamName = testData.teams[__ITER % testData.teams.length];
  const response = http.get(`${baseURL}/team/get?team_name=${teamName}`, {
    timeout: '2s'
  });
  
  check(response, {
    'get team status 200': (r) => r.status === 200,
  });
}

function getStats() {
  const response = http.get(`${baseURL}/stats/reviews`, {
    timeout: '2s'
  });
  
  check(response, {
    'get stats status 200': (r) => r.status === 200,
  });
}

function getUserReviews() {
  const activeUsers = testData.users.filter(u => u.is_active);
  if (activeUsers.length === 0) return;
  
  const user = activeUsers[__ITER % activeUsers.length];
  const response = http.get(`${baseURL}/users/getReview?user_id=${user.user_id}`, {
    timeout: '2s'
  });
  
  check(response, {
    'get user reviews valid status': (r) => r.status === 200 || r.status === 404,
  });
}

function getTeamAgain() {
  const teamName = testData.teams[(__ITER + 1) % testData.teams.length];
  const response = http.get(`${baseURL}/team/get?team_name=${teamName}`, {
    timeout: '2s'
  });
  
  check(response, {
    'get team again status 200': (r) => r.status === 200,
  });
}

function getStatsAgain() {
  const response = http.get(`${baseURL}/stats/reviews`, {
    timeout: '2s'
  });
  
  check(response, {
    'get stats again status 200': (r) => r.status === 200,
  });
}

function createPR() {
  const uniqueId = `${__VU}-${__ITER}-${Date.now()}-${Math.random().toString(36).substr(2, 8)}`;
  const prId = `pr-unique-${uniqueId}`;
  
  const activeUsers = testData.users.filter(u => u.is_active);
  if (activeUsers.length === 0) return;
  
  const author = activeUsers[__ITER % activeUsers.length];
  
  const payload = {
    pull_request_id: prId,
    pull_request_name: `Unique PR ${prId}`,
    author_id: author.user_id
  };
  
  const response = http.post(`${baseURL}/pullRequest/create`, JSON.stringify(payload), {
    headers: { 'Content-Type': 'application/json' },
    timeout: '2s'
  });
  
  check(response, {
    'create PR status 201': (r) => r.status === 201,
  });
}

export function teardown(data) {
}