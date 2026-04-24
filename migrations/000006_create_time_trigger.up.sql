CREATE FUNCTION fnc_trg_update_last_update_date()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_update_date = CURRENT_DATE;
    RETURN NEW;
END;
$$ language plpgsql;

CREATE OR REPLACE TRIGGER trg_update_products_last_update_date
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION fnc_trg_update_last_update_date();
