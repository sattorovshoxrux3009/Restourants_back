CREATE TABLE IF NOT EXISTS `restaurants` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `address` VARCHAR(255) NOT NULL,
    `latitude` DECIMAL(10, 7) NOT NULL,
    `longitude` DECIMAL(10, 7) NOT NULL,
    `phone_number` VARCHAR(20) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `capacity` INT NOT NULL,
    `owner_id` INT UNSIGNED NOT NULL,
    `opening_hours` VARCHAR(255) NOT NULL,
    `image_url` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `alcohol_permission` BOOLEAN NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `status` ENUM('active', 'inactive') NOT NULL DEFAULT 'active'
);
 
CREATE TABLE IF NOT EXISTS `admins` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `phone_number` VARCHAR(20) NOT NULL,
    `username` VARCHAR(100) NOT NULL UNIQUE,
    `password_hash` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `status` ENUM('active', 'inactive') NOT NULL DEFAULT 'inactive'
);

CREATE TABLE IF NOT EXISTS `event_prices` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `restaurant_id` INT UNSIGNED NOT NULL,
    `event_type` ENUM('morning', 'night') NOT NULL,
    `table_price` DECIMAL(8, 2) NOT NULL,
    `waiter_price` DECIMAL(8, 2) NOT NULL,
    `max_guests` INT NOT NULL,
    `table_seats` INT NOT NULL,
    `max_waiters` INT NOT NULL,
    `alcohol_permission` BOOLEAN NOT NULL,
    FOREIGN KEY (`restaurant_id`) REFERENCES `Restaurants`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `menu` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `restaurant_id` INT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `price` DECIMAL(8, 2) NOT NULL,
    `image_url` VARCHAR(255) NOT NULL,
    `category` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`restaurant_id`) REFERENCES `Restaurants`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `tokens` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `admin_id` INT UNSIGNED NOT NULL,
    `token` VARCHAR(255) NOT NULL,
    `auth_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`admin_id`) REFERENCES `Admins`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `super_admins` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name` VARCHAR(255) NOT NULL,
    `username` VARCHAR(100) NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `token` VARCHAR(255) ,
    `last_login` TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `admin_restaurant_limits` (
    `id`              INT AUTO_INCREMENT PRIMARY KEY,
    `admin_id`        INT UNSIGNED  NOT NULL,
    `max_restaurants` INT NOT NULL DEFAULT 1, 
    `created_at`      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE
);


ALTER TABLE `restaurants` ADD CONSTRAINT `restaurants_owner_id_fk` FOREIGN KEY (`owner_id`) REFERENCES `admins`(`id`) ON DELETE CASCADE;

