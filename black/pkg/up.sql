CREATE TABLE datasets (
  id BIGSERIAL PRIMARY KEY,
  key VARCHAR NOT NULL,  
  name VARCHAR NOT NULL,
  tags text[] NOT NULL,
  license TEXT NOT NULL,
  document_count BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT dataset_key UNIQUE (key)
);


CREATE TABLE documents (
  id BIGSERIAL PRIMARY KEY,
  dataset_key VARCHAR NOT NULL REFERENCES datasets(key) ON DELETE CASCADE,
  document_key VARCHAR NOT NULL,
  tags text[] NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  data BYTEA NOT NULL,
  CONSTRAINT dataset_key_document_key UNIQUE (dataset_key, document_key)
);


CREATE INDEX idx_datasets_tags ON datasets USING GIN(tags);
CREATE INDEX idx_documents_tags ON documents USING GIN(tags);