-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop tables if they exist (for clean setup)
DROP TABLE IF EXISTS blog_tags;
DROP TABLE IF EXISTS contact_submissions;
DROP TABLE IF EXISTS portfolio_projects;
DROP TABLE IF EXISTS blog_posts;
DROP TABLE IF EXISTS tags;

-- Create tables with proper timestamps and constraints

-- Tags table
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT tags_name_unique UNIQUE (name)
);

-- Blog posts table
CREATE TABLE blog_posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    summary TEXT,
    image_url VARCHAR(255),
    published BOOLEAN DEFAULT FALSE,
    publish_at TIMESTAMP WITH TIME ZONE,
    slug VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT blog_posts_slug_unique UNIQUE (slug)
);

-- Join table for blog posts and tags (many-to-many)
CREATE TABLE blog_tags (
    blog_post_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (blog_post_id, tag_id),
    FOREIGN KEY (blog_post_id) REFERENCES blog_posts(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Portfolio projects table
CREATE TABLE portfolio_projects (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    project_type VARCHAR(50) NOT NULL,
    image_url VARCHAR(255),
    project_url VARCHAR(255),
    repo_url VARCHAR(255),
    technologies JSONB,
    featured BOOLEAN DEFAULT FALSE,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Contact submissions table
CREATE TABLE contact_submissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    message TEXT NOT NULL,
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX idx_blog_posts_slug ON blog_posts(slug);
CREATE INDEX idx_blog_posts_published ON blog_posts(published);
CREATE INDEX idx_portfolio_projects_type ON portfolio_projects(project_type);
CREATE INDEX idx_portfolio_projects_featured ON portfolio_projects(featured);
CREATE INDEX idx_contact_submissions_read ON contact_submissions(read);

-- Create updated_at triggers
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to all tables
CREATE TRIGGER update_blog_posts_timestamp
BEFORE UPDATE ON blog_posts
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_tags_timestamp
BEFORE UPDATE ON tags
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_portfolio_projects_timestamp
BEFORE UPDATE ON portfolio_projects
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_contact_submissions_timestamp
BEFORE UPDATE ON contact_submissions
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();