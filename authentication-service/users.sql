-- 1. Sequence Setup
-- Create an auto-incrementing sequence for the user ID
CREATE SEQUENCE public.user_id_seq 
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;

ALTER TABLE public.user_id_seq OWNER TO postgres;

-- 2. Table Definition
-- Defines the 'users' table structure with 2-space indentation
CREATE TABLE public.users (
  id          integer      NOT NULL DEFAULT nextval('public.user_id_seq'::regclass),
  email       varchar(255),
  first_name  varchar(255),
  last_name   varchar(255),
  password    varchar(60),
  user_active integer      DEFAULT 0,
  created_at  timestamp    without time zone,
  updated_at  timestamp    without time zone,
  
  -- Set the primary key
  CONSTRAINT users_pkey PRIMARY KEY (id)
);

ALTER TABLE public.users OWNER TO postgres;

-- 3. Sequence Sync
-- Ensure the sequence starts at the correct value
SELECT pg_catalog.setval('public.user_id_seq', 1, true);

-- 4. Initial Data Seed
-- Insert the default administrator account
INSERT INTO public.users (
  email, 
  first_name, 
  last_name, 
  password, 
  user_active, 
  created_at, 
  updated_at
)
VALUES (
  'admin@example.com',
  'Admin',
  'System',
  '$2a$12$Gr2jYdUz7JRqL3v.b1xDnOMGhVBY.ATcLhjNkXozH4TSOhL1I1dTa',
  1,
  '2022-03-14 00:00:00',
  '2022-03-14 00:00:00'
);
