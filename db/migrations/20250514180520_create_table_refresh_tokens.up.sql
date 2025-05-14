CREATE TABLE refresh_tokens (
  admin_id INT NOT NULL,
  token TEXT NOT NULL,
  status_logout TINYINT NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (admin_id),
  FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE
);