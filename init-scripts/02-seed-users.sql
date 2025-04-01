INSERT INTO users (username, email, password_hash, role, active)
VALUES (
    'admin',
    'limanjaya.jason@gmail.com',
    '$2a$10$Am6mRAdXsxw9VyRU0GMhJ.V7CfLLgZlGvySslZB99SgGPhxyHaQwC',
    'admin',
    true
) ON CONFLICT (username) DO NOTHING;