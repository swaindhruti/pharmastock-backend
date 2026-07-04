CREATE TABLE inventories (
    stockist_id BIGINT NOT NULL REFERENCES stockists(id) ON DELETE CASCADE,
    medicine_id BIGINT NOT NULL REFERENCES medicines(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (stockist_id, medicine_id)
);

CREATE INDEX idx_inventories_medicine_id ON inventories (medicine_id);
