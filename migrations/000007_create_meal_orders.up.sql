CREATE TABLE IF NOT EXISTS meal_orders (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    user_id INT UNSIGNED NOT NULL,
    week_start_date DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE KEY unique_user_week (user_id, week_start_date)
);

CREATE TABLE IF NOT EXISTS meal_order_items (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    meal_order_id INT UNSIGNED NOT NULL,
    day_of_week ENUM('Monday', 'Tuesday', 'Wednesday', 'Thursday') NOT NULL,
    meal_id INT UNSIGNED NOT NULL,
    FOREIGN KEY (meal_order_id) REFERENCES meal_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (meal_id) REFERENCES meals(id),
    UNIQUE KEY unique_order_day (meal_order_id, day_of_week)
);
