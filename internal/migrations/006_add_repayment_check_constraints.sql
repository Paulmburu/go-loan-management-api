ALTER TABLE repayments
ADD CONSTRAINT repayments_amount_positive CHECK (amount > 0);