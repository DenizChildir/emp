#!/bin/bash

# Update package list
sudo apt update

# Install MySQL
sudo apt install -y mysql-server

# Start MySQL service
sudo systemctl start mysql

# Secure MySQL installation and set root password
sudo mysql_secure_installation <<EOF

y
2
uuuu
uuuu
y
y
y
y
EOF

# Create SQL file
cat > setup.sql <<EOL
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS \`mydb\` DEFAULT CHARACTER SET utf8mb3 ;
USE \`mydb\` ;

-- -----------------------------------------------------
-- Table \`mydb\`.\`customer\`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS \`mydb\`.\`customer\` (
  \`id\` INT NOT NULL AUTO_INCREMENT,
  \`fullName\` VARCHAR(255) NULL DEFAULT NULL,
  \`company\` VARCHAR(255) NULL DEFAULT NULL,
  \`managerEmail\` VARCHAR(255) NULL DEFAULT NULL,
  \`ccEmail\` VARCHAR(255) NULL DEFAULT NULL,
  \`dutyLunch\` TINYINT NULL DEFAULT NULL,
  PRIMARY KEY (\`id\`))
ENGINE = InnoDB
AUTO_INCREMENT = 6
DEFAULT CHARACTER SET = utf8mb3;

-- -----------------------------------------------------
-- Table \`mydb\`.\`employee\`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS \`mydb\`.\`employee\` (
  \`id\` INT NOT NULL AUTO_INCREMENT,
  \`fullName\` VARCHAR(255) NULL DEFAULT NULL,
  \`active\` TINYINT NULL DEFAULT NULL,
  \`phone\` VARCHAR(20) NULL DEFAULT NULL,
  \`email\` VARCHAR(255) NULL DEFAULT NULL,
  \`MsId\` INT NULL DEFAULT NULL,
  \`excalID\` INT NULL DEFAULT NULL,
  \`rate\` DECIMAL(10,2) NULL DEFAULT NULL,
  PRIMARY KEY (\`id\`))
ENGINE = InnoDB
AUTO_INCREMENT = 4
DEFAULT CHARACTER SET = utf8mb3;

-- -----------------------------------------------------
-- Table \`mydb\`.\`day\`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS \`mydb\`.\`day\` (
  \`id\` INT NOT NULL AUTO_INCREMENT,
  \`date\` DATE NULL DEFAULT NULL,
  \`startTime\` TIME NULL DEFAULT NULL,
  \`endTime\` TIME NULL DEFAULT NULL,
  \`employeeId\` INT NULL DEFAULT NULL,
  \`shiftHours\` DECIMAL(5,2) NULL DEFAULT NULL,
  \`customerId\` INT NULL DEFAULT NULL,
  \`startMealTime\` TIME NULL DEFAULT NULL,
  \`endMealTime\` TIME NULL DEFAULT NULL,
  PRIMARY KEY (\`id\`),
  INDEX \`employeeId\` (\`employeeId\` ASC) VISIBLE,
  INDEX \`customerId\` (\`customerId\` ASC) VISIBLE,
  CONSTRAINT \`day_ibfk_1\`
    FOREIGN KEY (\`employeeId\`)
    REFERENCES \`mydb\`.\`employee\` (\`id\`),
  CONSTRAINT \`day_ibfk_2\`
    FOREIGN KEY (\`customerId\`)
    REFERENCES \`mydb\`.\`customer\` (\`id\`))
ENGINE = InnoDB
AUTO_INCREMENT = 22
DEFAULT CHARACTER SET = utf8mb3;

-- -----------------------------------------------------
-- Table \`mydb\`.\`weektotal\`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS \`mydb\`.\`weektotal\` (
  \`id\` INT NOT NULL AUTO_INCREMENT,
  \`employeeId\` INT NULL DEFAULT NULL,
  \`weekStartDate\` DATE NULL DEFAULT NULL,
  \`sundayId\` INT NULL DEFAULT NULL,
  \`mondayId\` INT NULL DEFAULT NULL,
  \`tuesdayId\` INT NULL DEFAULT NULL,
  \`wednesdayId\` INT NULL DEFAULT NULL,
  \`thursdayId\` INT NULL DEFAULT NULL,
  \`fridayId\` INT NULL DEFAULT NULL,
  \`saturdayId\` INT NULL DEFAULT NULL,
  \`sundayHours\` DECIMAL(5,2) NULL DEFAULT NULL,
  \`mondayHours\` DECIMAL(5,2) NULL DEFAULT NULL,
  \`tuesdayHours\` DECIMAL(5,2) NULL DEFAULT NULL,
  \`wednesdayHours\` DECIMAL(5,2) NULL DEFAULT NULL,
  \`thursdayHours\` DECIMAL(5,2) NULL DEFAULT NULL,
  \`fridayHours\` DECIMAL(5,2) NULL DEFAULT NULL,
  \`saturdayHours\` DECIMAL(5,2) NULL DEFAULT NULL,
  \`finalHours\` DECIMAL(5,2) NULL DEFAULT NULL,
  PRIMARY KEY (\`id\`),
  INDEX \`employeeId\` (\`employeeId\` ASC) VISIBLE,
  INDEX \`sundayId\` (\`sundayId\` ASC) VISIBLE,
  INDEX \`mondayId\` (\`mondayId\` ASC) VISIBLE,
  INDEX \`tuesdayId\` (\`tuesdayId\` ASC) VISIBLE,
  INDEX \`wednesdayId\` (\`wednesdayId\` ASC) VISIBLE,
  INDEX \`thursdayId\` (\`thursdayId\` ASC) VISIBLE,
  INDEX \`fridayId\` (\`fridayId\` ASC) VISIBLE,
  INDEX \`saturdayId\` (\`saturdayId\` ASC) VISIBLE,
  CONSTRAINT \`weektotal_ibfk_1\`
    FOREIGN KEY (\`employeeId\`)
    REFERENCES \`mydb\`.\`employee\` (\`id\`),
  CONSTRAINT \`weektotal_ibfk_2\`
    FOREIGN KEY (\`sundayId\`)
    REFERENCES \`mydb\`.\`day\` (\`id\`),
  CONSTRAINT \`weektotal_ibfk_3\`
    FOREIGN KEY (\`mondayId\`)
    REFERENCES \`mydb\`.\`day\` (\`id\`),
  CONSTRAINT \`weektotal_ibfk_4\`
    FOREIGN KEY (\`tuesdayId\`)
    REFERENCES \`mydb\`.\`day\` (\`id\`),
  CONSTRAINT \`weektotal_ibfk_5\`
    FOREIGN KEY (\`wednesdayId\`)
    REFERENCES \`mydb\`.\`day\` (\`id\`),
  CONSTRAINT \`weektotal_ibfk_6\`
    FOREIGN KEY (\`thursdayId\`)
    REFERENCES \`mydb\`.\`day\` (\`id\`),
  CONSTRAINT \`weektotal_ibfk_7\`
    FOREIGN KEY (\`fridayId\`)
    REFERENCES \`mydb\`.\`day\` (\`id\`),
  CONSTRAINT \`weektotal_ibfk_8\`
    FOREIGN KEY (\`saturdayId\`)
    REFERENCES \`mydb\`.\`day\` (\`id\`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
EOL

# Run SQL script
sudo mysql -u root -puuuu < setup.sql

echo "MySQL installation and database setup complete."