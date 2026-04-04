CREATE TABLE IF NOT EXISTS meal_orders (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    user_id INT UNSIGNED NOT NULL,
    week_start_date DATE NOT NULL,
    monday_meal_id INT UNSIGNED NULL,
    tuesday_meal_id INT UNSIGNED NULL,
    wednesday_meal_id INT UNSIGNED NULL,
    thursday_meal_id INT UNSIGNED NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (monday_meal_id) REFERENCES meals(id),
    FOREIGN KEY (tuesday_meal_id) REFERENCES meals(id),
    FOREIGN KEY (wednesday_meal_id) REFERENCES meals(id),
    FOREIGN KEY (thursday_meal_id) REFERENCES meals(id),
    UNIQUE KEY unique_user_week (user_id, week_start_date)
);
