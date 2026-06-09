CREATE TABLE jobs (
    id BIGSERIAL PRIMARY KEY,
    stockist_id BIGINT NOT NULL,
    job_status VARCHAR(20) NOT NULL CHECK (job_status IN ('pending', 'processing', 'completed', 'failed')),
    file_path TEXT NOT NULL,
    error_message TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,

    FOREIGN KEY (stockist_id) REFERENCES stockists(id) ON DELETE CASCADE
);

CREATE INDEX idx_jobs_status ON jobs(job_status);
CREATE INDEX idx_jobs_stockist_id ON jobs(stockist_id);