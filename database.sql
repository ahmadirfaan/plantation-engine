DO
$$
    BEGIN
        IF NOT EXISTS (SELECT
                       FROM pg_database
                       WHERE datname = 'database') THEN
            PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE database');
        END IF;
    END
$$;

--connect the database
\c database

CREATE EXTENSION IF NOT EXISTS dblink;


CREATE TABLE IF NOT EXISTS estate
(
    id         text PRIMARY KEY,
    name       VARCHAR(255),
    width      INT,
    length     INT,
    ext_info   TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS block
(
    id           UUID,
    estate_id    UUID NOT NULL,
    x_coordinate INT,
    y_coordinate INT,
    ext_info     TEXT,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id, estate_id)

) PARTITION BY HASH (estate_id);

DO
$$
    DECLARE
        i         INT;
        part_name TEXT;
    BEGIN
        FOR i IN 0..99
            LOOP
                part_name := format('block_p%s', i);

                -- Create partition if it doesn't exist
                IF NOT EXISTS (SELECT
                               FROM pg_class
                               WHERE relname = part_name) THEN
                    EXECUTE format('
                CREATE TABLE %I PARTITION OF block
                FOR VALUES WITH (MODULUS 100, REMAINDER %s);
            ', part_name, i);
                END IF;

                -- Create unique index if it doesn't exist
                IF NOT EXISTS (SELECT
                               FROM pg_indexes
                               WHERE tablename = part_name AND indexname = part_name || '_uniq_estate_xy') THEN
                    EXECUTE format('
                CREATE UNIQUE INDEX %I ON %I (estate_id, x_coordinate, y_coordinate);
            ', part_name || '_uniq_estate_xy', part_name);
                END IF;

            END LOOP;
    END
$$;


CREATE TABLE IF NOT EXISTS tree
(
    id         UUID,
    block_id   UUID NOT NULL,
    estate_id  UUID NOT NULL,
    height     INT,
    ext_info   TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id, estate_id)
) PARTITION BY HASH (estate_id);

DO
$$
    DECLARE
        i          INT;
        part_name  TEXT;
        index_name TEXT;
    BEGIN
        FOR i IN 0..99
            LOOP
                part_name := format('tree_p%s', i);
                index_name := part_name || '_uniq_block_estate';

                -- Create partition if not exists
                IF NOT EXISTS (SELECT
                               FROM pg_class
                               WHERE relname = part_name) THEN
                    EXECUTE format('
                CREATE TABLE %I PARTITION OF tree
                FOR VALUES WITH (MODULUS 100, REMAINDER %s);
            ', part_name, i);
                END IF;

                -- Create unique index if not exists
                IF NOT EXISTS (SELECT FROM pg_indexes WHERE tablename = part_name AND indexname = index_name) THEN
                    EXECUTE format('
                CREATE UNIQUE INDEX %I ON %I (block_id, estate_id);
            ', index_name, part_name);
                END IF;

            END LOOP;
    END
$$;


CREATE TABLE IF NOT EXISTS estate_stats
(
    estate_id            UUID PRIMARY KEY,
    min_height_tree      DOUBLE PRECISION,
    max_height_tree      DOUBLE PRECISION,
    sum_tree             INT,
    median_height_tree   DOUBLE PRECISION,
    total_distance_drone DOUBLE PRECISION,
    ext_info             TEXT,
    created_at           TIMESTAMP DEFAULT NOW(),
    updated_at           TIMESTAMP DEFAULT NOW(),
    calculated_at        TIMESTAMP DEFAULT NOW()
);
