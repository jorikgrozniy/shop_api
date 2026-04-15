CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_name VARCHAR(100) NOT NULL,
    client_surname VARCHAR(100) NOT NULL,
    birthdate DATE NOT NULL,
    gender VARCHAR CHECK (gender IN ('male', 'female')),
    registration_date DATE DEFAULT CURRENT_DATE,
    address_id UUID NOT NULL,

    CONSTRAINT fk_client_address FOREIGN KEY (address_id) REFERENCES addresses(id)
);