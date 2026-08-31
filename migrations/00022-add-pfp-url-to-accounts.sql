ALTER TABLE accounts
  ADD COLUMN pfp_url VARCHAR(255) NOT NULL DEFAULT '/assets/default_pfp.png';
