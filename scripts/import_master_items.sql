-- Import master items from the PROGAS CSV.
-- Run with: mysql --local-infile=1 -h HOST -P PORT -u USER -p DATABASE < scripts/import_master_items.sql

DROP TEMPORARY TABLE IF EXISTS _master_items_import;

CREATE TEMPORARY TABLE _master_items_import (
    id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    sku VARCHAR(50) NOT NULL,
    gas_type VARCHAR(50),
    is_serialized VARCHAR(10) NOT NULL,
    empty_weight_kg VARCHAR(32),
    gas_weight_kg VARCHAR(32),
    min_stock_alert VARCHAR(32),
    satuan VARCHAR(50),
    kategori VARCHAR(128),
    PRIMARY KEY (id),
    UNIQUE KEY uq_import_sku (sku)
) ENGINE = InnoDB;

LOAD DATA LOCAL INFILE '/Users/debi.darmawan_aam/work/github.com/progas/master_items_progasindonesia.csv'
INTO TABLE _master_items_import
FIELDS TERMINATED BY ','
OPTIONALLY ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 LINES
(@name, @sku, @gas_type, @is_serialized, @empty_weight_kg,
 @gas_weight_kg, @min_stock_alert, @satuan, @kategori)
SET
    id = UUID(),
    name = TRIM(@name),
    sku = TRIM(@sku),
    gas_type = TRIM(@gas_type),
    is_serialized = LOWER(TRIM(@is_serialized)),
    empty_weight_kg = TRIM(@empty_weight_kg),
    gas_weight_kg = TRIM(@gas_weight_kg),
    min_stock_alert = TRIM(@min_stock_alert),
    satuan = TRIM(@satuan),
    kategori = TRIM(@kategori);

-- Abort before writing if the CSV violates the application rules.
DELIMITER //
CREATE PROCEDURE _validate_master_items_import()
BEGIN
    IF EXISTS (
        SELECT 1
        FROM _master_items_import
        WHERE name = ''
           OR sku = ''
           OR is_serialized NOT IN ('true', 'false')
           OR (is_serialized = 'true' AND COALESCE(gas_type, '') = '')
           OR (is_serialized = 'false' AND COALESCE(gas_type, '') <> '')
           OR (empty_weight_kg <> '' AND empty_weight_kg NOT REGEXP '^[0-9]+(\\.[0-9]+)?$')
           OR (gas_weight_kg <> '' AND gas_weight_kg NOT REGEXP '^[0-9]+(\\.[0-9]+)?$')
           OR (min_stock_alert <> '' AND min_stock_alert NOT REGEXP '^[0-9]+$')
    ) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'CSV contains invalid master item values';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM _master_items_import imported
        JOIN master_item existing
          ON existing.sku = imported.sku
         AND existing.deleted_at IS NULL
    ) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'One or more SKUs already exist';
    END IF;
END//
DELIMITER ;

CALL _validate_master_items_import();
DROP PROCEDURE _validate_master_items_import;

START TRANSACTION;

INSERT INTO master_item (
    id,
    name,
    sku,
    gas_type,
    is_serialized,
    empty_weight_kg,
    gas_weight_kg,
    min_stock_alert,
    max_days_at_customer,
    created_at,
    updated_at
)
SELECT
    id,
    name,
    sku,
    NULLIF(gas_type, ''),
    is_serialized = 'true',
    COALESCE(NULLIF(empty_weight_kg, ''), 0),
    COALESCE(NULLIF(gas_weight_kg, ''), 0),
    COALESCE(NULLIF(min_stock_alert, ''), 10),
    0,
    NOW(),
    NOW()
FROM _master_items_import;

INSERT INTO sparepart_stock (
    id,
    item_id,
    quantity,
    created_at,
    updated_at
)
SELECT
    UUID(),
    id,
    0,
    NOW(),
    NOW()
FROM _master_items_import
WHERE is_serialized = 'false';

COMMIT;

SELECT COUNT(*) AS imported_master_items
FROM master_item
WHERE sku IN (SELECT sku FROM _master_items_import)
  AND deleted_at IS NULL;

SELECT COUNT(*) AS imported_non_serialized_stock_rows
FROM sparepart_stock
WHERE item_id IN (SELECT id FROM _master_items_import);
