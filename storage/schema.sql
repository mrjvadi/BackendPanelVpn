
-- ==========================================
-- vpn-panel final schema - ترکیبی بهینه‌شده
-- ==========================================

-- جداول پایه
CREATE TABLE migrations (
                            id SERIAL PRIMARY KEY,
                            name TEXT,
                            applied_at TIMESTAMP
);

CREATE TABLE admins (
                        id SERIAL PRIMARY KEY,
                        username VARCHAR(100) UNIQUE,
                        password VARCHAR(255),
                        telegram_id BIGINT UNIQUE,
                        email VARCHAR(100),
                        created_at TIMESTAMP,
                        updated_at TIMESTAMP
);

CREATE TABLE servers (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(100),
                         version VARCHAR(50),
                         url VARCHAR(255),
                         unique_name varchar(100),
                         username VARCHAR(100),
                         password VARCHAR(255),
                         remove_prefix VARCHAR(100),
                         total_quota BIGINT,
                         used_quota BIGINT,
                         remaining_quota BIGINT,
                         created_at TIMESTAMP,
                         updated_at TIMESTAMP
);

CREATE TABLE configs (
                         id SERIAL PRIMARY KEY,
                         server_id INT REFERENCES servers(id),
                         link TEXT UNIQUE,
                         tag VARCHAR(100),
                         custom_tag VARCHAR(100),
                         is_active BOOLEAN,
                         created_at TIMESTAMP,
                         updated_at TIMESTAMP
);

CREATE TABLE resellers (
                           id SERIAL PRIMARY KEY,
                           parent_reseller_id INT NULL REFERENCES resellers(id),
                           username VARCHAR(100) UNIQUE,
                           password VARCHAR(255),
                           telegram_id BIGINT UNIQUE,
                           email VARCHAR(100),
                           api_token VARCHAR(512),
                           custom_name VARCHAR(100) UNIQUE,
                           total_quota BIGINT,
                           used_quota BIGINT,
                           remaining_quota BIGINT,
                           status VARCHAR(20),
                           can_create_subreseller BOOLEAN,
                           created_at TIMESTAMP,
                           updated_at TIMESTAMP
);

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       reseller_id INT REFERENCES resellers(id),
                       username VARCHAR(100) UNIQUE,
                       uuid UUID UNIQUE,
                       total_quota BIGINT,
                       used_quota BIGINT,
                       remaining_quota BIGINT,
                       activation_date TIMESTAMP,
                       status VARCHAR(20),
                       sub_token VARCHAR(512),
                       time_limit TIMESTAMP,
                       created_at TIMESTAMP,
                       updated_at TIMESTAMP
);

-- تگ سیستم جدید
CREATE TABLE tag_system (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(100) UNIQUE,
                            is_active BOOLEAN,

                            created_at TIMESTAMP,
                            updated_at TIMESTAMP
);

CREATE TABLE tag_system_config (
                                   id SERIAL PRIMARY KEY,
                                   tag_system_id INT REFERENCES tag_system(id),
                                   config_id INT REFERENCES configs(id)
);

CREATE TABLE reseller_tag_system_allowed (
                                             id SERIAL PRIMARY KEY,
                                             tag_system_id INT REFERENCES tag_system(id),
                                             reseller_id INT REFERENCES resellers(id)
);

CREATE TABLE tag_reseller (
                              id SERIAL PRIMARY KEY,
                              name VARCHAR(100) UNIQUE,
                              is_active BOOLEAN
);

CREATE TABLE tag_reseller_config (
                                     id SERIAL PRIMARY KEY,
                                     tag_reseller_id INT REFERENCES tag_reseller(id),
                                     reseller_id INT REFERENCES resellers(id),
                                     config_id INT REFERENCES configs(id),
                                     config_custom_name VARCHAR(100)
);

CREATE TABLE user_tag_reseller (
                                   id SERIAL PRIMARY KEY,
                                   user_id INT REFERENCES users(id),
                                   tag_reseller_id INT REFERENCES tag_reseller(id)
);

-- ادامه جداول اصلی
CREATE TABLE reseller_config (
                                 reseller_id INT REFERENCES resellers(id),
                                 config_id INT REFERENCES configs(id),
                                 is_active BOOLEAN,
                                 custom_name VARCHAR(100),
                                 created_at TIMESTAMP,
                                 updated_at TIMESTAMP,
                                 PRIMARY KEY (reseller_id, config_id)
);

CREATE TABLE user_config (
                             user_id INT REFERENCES users(id),
                             config_id INT REFERENCES configs(id),
                             assigned_at TIMESTAMP,
                             PRIMARY KEY (user_id, config_id)
);

CREATE TABLE usage_logs (
                            id SERIAL PRIMARY KEY,
                            user_id INT REFERENCES users(id),
                            reseller_id INT REFERENCES resellers(id),
                            server_id INT REFERENCES servers(id),
                            amount BIGINT,
                            created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_proxies (
                              id SERIAL PRIMARY KEY,
                              user_id INT REFERENCES users(id),
                              server_id INT REFERENCES servers(id),
                              type VARCHAR(50),
                              uuid UUID,
                              payload JSONB,
                              created_at TIMESTAMP,
                              UNIQUE (server_id, uuid)
);

CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              reseller_id INT REFERENCES resellers(id),
                              type VARCHAR(20),
                              amount BIGINT,
                              performed_by INT REFERENCES resellers(id),
                              description TEXT,
                              reference_type VARCHAR(50),
                              reference_id INT,
                              created_at TIMESTAMP
);

CREATE TABLE filtered_inbounds (
                                   id SERIAL PRIMARY KEY,
                                   inbound_tag VARCHAR(100) UNIQUE
);

CREATE TABLE filtered_tags (
                               id SERIAL PRIMARY KEY,
                               tag VARCHAR(100) UNIQUE
);

CREATE TABLE activity_logs (
                               id SERIAL PRIMARY KEY,
                               event_type VARCHAR(100) NOT NULL,
                               entity VARCHAR(100) NOT NULL,
                               entity_id INT NOT NULL,
                               request_id VARCHAR(255),
                               performed_by VARCHAR(100),
                               ip_address VARCHAR(45),
                               user_agent VARCHAR(255),
                               created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
                               updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
                               deleted_at TIMESTAMP WITHOUT TIME ZONE NULL
);


-- تابع بروزرسانی کوتاژ
CREATE OR REPLACE FUNCTION fn_update_quota()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE users SET used_quota = used_quota + NEW.amount WHERE id = NEW.user_id;
    WITH RECURSIVE chain AS (
        SELECT id, parent_reseller_id FROM resellers WHERE id = NEW.reseller_id
        UNION ALL
        SELECT r.id, r.parent_reseller_id FROM resellers r JOIN chain c ON r.id = c.parent_reseller_id
    )
    UPDATE resellers SET used_quota = used_quota + NEW.amount FROM chain WHERE resellers.id = chain.id;
    UPDATE servers SET used_quota = used_quota + NEW.amount WHERE id = NEW.server_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- تریگر برای اعمال تابع بعد از درج لاگ مصرف
CREATE TRIGGER trg_usage_logs_after_insert
    AFTER INSERT ON usage_logs
    FOR EACH ROW
EXECUTE FUNCTION fn_update_quota();

-- جداول آرشیو مصرف
CREATE TABLE usage_logs_hourly (LIKE usage_logs INCLUDING ALL);
CREATE TABLE usage_logs_daily (LIKE usage_logs INCLUDING ALL);
CREATE TABLE usage_logs_weekly (LIKE usage_logs INCLUDING ALL);
CREATE TABLE usage_logs_monthly (LIKE usage_logs INCLUDING ALL);

-- توابع آرشیو کردن
CREATE OR REPLACE FUNCTION archive_usage_logs_hourly() RETURNS void AS $$
BEGIN
    INSERT INTO usage_logs_hourly SELECT * FROM usage_logs WHERE created_at < NOW() - INTERVAL '1 hour';
    DELETE FROM usage_logs WHERE created_at < NOW() - INTERVAL '1 hour';
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION archive_usage_logs_daily() RETURNS void AS $$
BEGIN
    INSERT INTO usage_logs_daily SELECT * FROM usage_logs WHERE created_at < NOW() - INTERVAL '1 day';
    DELETE FROM usage_logs WHERE created_at < NOW() - INTERVAL '1 day';
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION archive_usage_logs_weekly() RETURNS void AS $$
BEGIN
    INSERT INTO usage_logs_weekly SELECT * FROM usage_logs WHERE created_at < NOW() - INTERVAL '1 week';
    DELETE FROM usage_logs WHERE created_at < NOW() - INTERVAL '1 week';
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION archive_usage_logs_monthly() RETURNS void AS $$
BEGIN
    INSERT INTO usage_logs_monthly SELECT * FROM usage_logs WHERE created_at < NOW() - INTERVAL '1 month';
    DELETE FROM usage_logs WHERE created_at < NOW() - INTERVAL '1 month';
END;
$$ LANGUAGE plpgsql;