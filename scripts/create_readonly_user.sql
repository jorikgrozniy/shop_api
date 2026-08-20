SET app.ro_user = :'ro_user';
SET app.ro_password = :'ro_password';

DO $$
DECLARE
    v_ro_user text := current_setting('app.ro_user');
    v_ro_password text := current_setting('app.ro_password');
BEGIN
    EXECUTE format(
        'CREATE ROLE %I LOGIN PASSWORD %L',
        v_ro_user,
        v_ro_password
    );

    EXECUTE format(
        'GRANT CONNECT ON DATABASE %I TO %I',
        current_database(),
        v_ro_user
    );

    EXECUTE format(
        'GRANT USAGE ON SCHEMA public TO %I',
        v_ro_user
    );

    EXECUTE format(
        'GRANT SELECT ON ALL TABLES IN SCHEMA public TO %I',
        v_ro_user
    );

    EXECUTE format(
        'ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO %I',
        v_ro_user
    );
END
$$;