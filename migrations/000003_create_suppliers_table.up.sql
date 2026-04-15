CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    address_id UUID NOT NULL,
    phone_number VARCHAR(20) NOT NULL,

    CONSTRAINT fk_supplier_address FOREIGN KEY (address_id) REFERENCES addresses(id)
);