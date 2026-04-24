CREATE TABLE IF NOT EXISTS repayments (
    id SERIAL PRIMARY KEY,
    loan_id INT NOT NULL REFERENCES loans(id),
    amount DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_repayments_loan_id ON repayments(loan_id);