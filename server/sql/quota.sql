
CREATE DATABASE IF NOT EXISTS dora_server ;

CREATE TABLE IF NOT EXISTS `quota` (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `username` VARCHAR(20) NOT NULL,
    `read_quota` INT NOT NULL DEFAULT 100000,
    `write_quota` INT NOT NULL DEFAULT 100000,
    PRIMARY KEY ( `id` ),
    UNIQUE (`username`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP PROCEDURE IF EXISTS proc_init_quota;
DELIMITER $
CREATE PROCEDURE proc_init_quota()
BEGIN
    DECLARE n INT DEFAULT 1 ;
    DECLARE rq INT DEFAULT 100000;
    DECLARE wq INT DEFAULT 0;
    DECLARE tmp VARCHAR(20) DEFAULT '';

    SET @execSql='insert into quota (username,read_quota, write_quota)  values';
    SET @execdata = '';

    WHILE n <= 1000001 DO
        SET rq = 100000;
        SET wq = 100000;
        SET tmp = CAST(n AS CHAR(20));
        IF n*3+3 < rq THEN SET rq = n*3+3; END IF;
        IF n+3 < wq THEN SET wq = n+3; END IF;

        SET @execdata=concat(@execdata,'(',tmp, ',', rq,',', wq,')');

        IF n%5000=0
        THEN
            set @execSql = concat(@execSql,@execdata,';');

            PREPARE stmt FROM @execSql;
            EXECUTE stmt;
            DEALLOCATE PREPARE stmt;
            COMMIT;

            SET @execSql='insert into quota (username,read_quota, write_quota) values';
            SET @execdata = '';
        ELSE
            SET @execdata = concat(@execdata,',');
        END IF;

        SET n = n + 1 ;
    END WHILE;
END $

CALL proc_init_quota();

