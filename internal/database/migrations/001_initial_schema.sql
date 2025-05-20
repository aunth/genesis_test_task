

CREATE TABLE IF NOT EXISTS subscription (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    city VARCHAR(100) NOT NULL,
    frequency VARCHAR(20) NOT NULL,
    confirmed BOOLEAN DEFAULT FALSE,
    confirmation_token VARCHAR(36) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subscription_email ON subscription(email);
CREATE INDEX IF NOT EXISTS idx_subscription_token ON subscription(confirmation_token);

CREATE INDEX IF NOT EXISTS idx_weather_city ON subscription(city); 