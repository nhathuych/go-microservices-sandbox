CREATE TABLE public.users (
  id          SERIAL PRIMARY KEY,
  email       varchar(255) NOT NULL UNIQUE,
  first_name  varchar(255),
  last_name   varchar(255),
  password    varchar(60)  NOT NULL,
  user_active integer      DEFAULT 0,
  created_at  timestamp    DEFAULT now(),
  updated_at  timestamp    DEFAULT now()
);

INSERT INTO public.users (email, first_name, last_name, password, user_active)
VALUES (
  'admin@example.com',
  'Admin',
  'System',
  '$2a$12$Gr2jYdUz7JRqL3v.b1xDnOMGhVBY.ATcLhjNkXozH4TSOhL1I1dTa',
  1
);
