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
    DECLARE tmp VARCHAR(20) DEFAULT '';
    DECLARE i INT DEFAULT 0;
    WHILE i<=1000000 DO
        SET tmp = CAST(i AS CHAR(20));
        INSERT INTO account(username, password) VALUES(tmp, tmp);
        SET i = i+1;
END WHILE;
END $
CALL proc_init_account();

