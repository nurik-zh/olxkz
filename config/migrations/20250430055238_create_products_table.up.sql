CREATE TABLE products (
      id SERIAL PRIMARY KEY,
      name TEXT NOT NULL,
      description TEXT,
      price DOUBLE PRECISION NOT NULL,
      category_id INTEGER NOT NULL
);
