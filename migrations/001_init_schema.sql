CREATE TABLE teams (
    team_id SERIAL PRIMARY KEY, 
    team_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE users (
    user_id VARCHAR(100) PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL,
    team_id INTEGER NOT NULL,
    FOREIGN KEY (team_id) REFERENCES teams(team_id) ON DELETE CASCADE
);

CREATE TABLE pull_requests (
    pull_request_id VARCHAR(100) PRIMARY KEY,
    pull_request_name VARCHAR(500) NOT NULL,
    author_id VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NULL,
    merged_at TIMESTAMP NULL,
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE pr_reviewers (
    pull_request_id VARCHAR(100) NOT NULL,
    reviewer_id VARCHAR(100) NOT NULL,
    PRIMARY KEY (pull_request_id, reviewer_id),
    FOREIGN KEY (pull_request_id) REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    FOREIGN KEY (reviewer_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_users_team_active ON users(team_id, is_active); 
CREATE INDEX idx_pr_author_status ON pull_requests(author_id, status);
CREATE INDEX idx_pr_reviewers_reviewer ON pr_reviewers(reviewer_id);
CREATE INDEX idx_pr_status ON pull_requests(status);
CREATE UNIQUE INDEX idx_teams_name ON teams(team_name); 