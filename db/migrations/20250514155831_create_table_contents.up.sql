CREATE TABLE contents (
  id INT AUTO_INCREMENT,
  title VARCHAR(150) NOT NULL,
  description TEXT NOT NULL,
  image TEXT,
  address TEXT,
  contact_info VARCHAR(100),
  category ENUM('kuliner', 'wisata', 'kerajinan') NOT NULL,
  created_by INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (created_by) REFERENCES admins(id) ON DELETE CASCADE
) ENGINE = InnoDB;