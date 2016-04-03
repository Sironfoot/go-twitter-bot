DROP DATABASE IF EXISTS go_twitter_bot;

CREATE DATABASE go_twitter_bot;
\connect go_twitter_bot;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
