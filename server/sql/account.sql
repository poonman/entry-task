CREATE DATABASE IF NOT EXISTS dora_server ;

CREATE TABLE IF NOT EXISTS `account` (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `username` VARCHAR(20) NOT NULL,
    `password` VARCHAR(20) NOT NULL,
    PRIMARY KEY ( `id` ),
    UNIQUE (`username`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP PROCEDURE IF EXISTS proc_init_account;
DELIMITER $
CREATE PROCEDURE proc_init_account()
BEGIN
    DECLARE n INT DEFAULT 1 ;
    DECLARE tmp VARCHAR(20) DEFAULT '';

    SET @execSql='insert into account (username,password)  values';
    SET @execdata = '';

    WHILE n <= 1000001 DO
        SET tmp = CAST(n AS CHAR(20));
        SET @execdata=concat(@execdata,'(',tmp,',', tmp,')');

        IF n%1000=0
        THEN
            set @execSql = concat(@execSql,@execdata,';');

            PREPARE stmt FROM @execSql;
            EXECUTE stmt;
            DEALLOCATE PREPARE stmt;
            COMMIT;

            set @execSql='insert into account (username,password)  values';
            set @execdata = '';
        ELSE
            set @execdata = concat(@execdata,',');
        end IF;

        SET n = n + 1 ;
    END WHILE ;
END $

CALL proc_init_account();

