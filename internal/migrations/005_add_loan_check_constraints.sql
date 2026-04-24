ALTER TABLE loans
ADD CONSTRAINT loans_principal_amount_positive CHECK (principal_amount > 0),
ADD CONSTRAINT loans_interest_rate_non_negative CHECK (interest_rate >= 0),
ADD CONSTRAINT loans_total_amount_non_negative CHECK (total_amount >= 0),
ADD CONSTRAINT loans_outstanding_amount_non_negative CHECK (outstanding_amount >= 0),
ADD CONSTRAINT loans_status_valid CHECK (status IN ('active', 'paid'));