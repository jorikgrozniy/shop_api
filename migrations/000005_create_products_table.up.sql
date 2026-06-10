CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    available_stock INT NOT NULL CHECK (available_stock >= 0),
    last_update_date DATE DEFAULT CURRENT_DATE,
    supplier_id UUID,
    image_id UUID,

    CONSTRAINT fk_product_supplier 
        FOREIGN KEY (supplier_id) REFERENCES suppliers(id)
        ON DELETE SET NULL,
    
    CONSTRAINT fk_product_image 
        FOREIGN KEY (image_id) REFERENCES images(id)
        ON DELETE SET NULL
);