-- Table: public.users

-- DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    cartid uuid NOT NULL DEFAULT uuid_generate_v4(),
    firstname character varying COLLATE pg_catalog."default" NOT NULL,
    lastname character varying COLLATE pg_catalog."default" NOT NULL,
    email character varying COLLATE pg_catalog."default" NOT NULL,
    password character varying COLLATE pg_catalog."default",
    isadmin boolean NOT NULL DEFAULT false,
    CONSTRAINT users_pkey PRIMARY KEY (email)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;


-- Table: public.carts

-- DROP TABLE IF EXISTS public.carts;

CREATE TABLE IF NOT EXISTS public.carts
(
    id uuid NOT NULL,
    totalitems integer NOT NULL DEFAULT 0,
    products jsonb,
    CONSTRAINT carts_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.carts
    OWNER to postgres;

-- Table: public.products

-- DROP TABLE IF EXISTS public.products;

CREATE TABLE IF NOT EXISTS public.products
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    price double precision NOT NULL,
    description character varying COLLATE pg_catalog."default" NOT NULL,
    image uuid NOT NULL DEFAULT uuid_generate_v4(),
    CONSTRAINT products_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.products
    OWNER to postgres;

