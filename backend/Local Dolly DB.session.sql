CREATE TABLE IF NOT EXISTS store(
  id UUID NOT NULL PRIMARY KEY,
  name VARCHAR(80) NOT NULL,
  email VARCHAR(80) NOT NULL,
  password VARCHAR(80)
);

INSERT INTO store(
  id,
  name,
  email,
  password)
VALUES ('e78a034e-81ec-498d-85c9-05fb9876c14c', 'Atrati', 'teste@test.com', 'admin123')