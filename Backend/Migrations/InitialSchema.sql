-- Extension: "uuid-ossp"

-- DROP EXTENSION "uuid-ossp";

CREATE EXTENSION "uuid-ossp"
    SCHEMA public
    VERSION "1.1";


-- Table: public.recipe

-- DROP TABLE public.recipe;

CREATE TABLE public.recipe
(
    id uuid NOT NULL DEFAULT uuid_generate_v1(),
    "time" interval NOT NULL,
    type text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT recipe_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.recipe
    OWNER to postgres;


-- Table: public.picture

-- DROP TABLE public.picture;

CREATE TABLE public.picture
(
    id uuid NOT NULL DEFAULT uuid_generate_v1(),
    recipe_id uuid NOT NULL,
    image_source text COLLATE pg_catalog."default" NOT NULL,
    sort_order numeric NOT NULL,
    CONSTRAINT picture_pkey PRIMARY KEY (id),
    CONSTRAINT picture_recipe_id_fkey FOREIGN KEY (recipe_id)
        REFERENCES public.recipe (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.picture
    OWNER to postgres;


-- Table: public.ingredient

-- DROP TABLE public.ingredient;

CREATE TABLE public.ingredient
(
    id uuid NOT NULL DEFAULT uuid_generate_v1(),
    recipe_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    quantity numeric NOT NULL,
    CONSTRAINT ingredient_pkey PRIMARY KEY (id),
    CONSTRAINT ingredient_recipe_id_fkey FOREIGN KEY (recipe_id)
        REFERENCES public.recipe (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.ingredient
    OWNER to postgres;


-- Table: public.equipment

-- DROP TABLE public.equipment;

CREATE TABLE public.equipment
(
    id uuid NOT NULL DEFAULT uuid_generate_v1(),
    recipe_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    quantity numeric NOT NULL,
    CONSTRAINT equipment_pkey PRIMARY KEY (id),
    CONSTRAINT equipment_recipe_id_fkey FOREIGN KEY (recipe_id)
        REFERENCES public.recipe (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.equipment
    OWNER to postgres;


-- Table: public.step

-- DROP TABLE public.step;

CREATE TABLE public.step
(
    id uuid NOT NULL DEFAULT uuid_generate_v1(),
    recipe_id uuid NOT NULL,
    content text COLLATE pg_catalog."default" NOT NULL,
    step_number numeric NOT NULL,
    CONSTRAINT step_pkey PRIMARY KEY (id),
    CONSTRAINT step_recipe_id_fkey FOREIGN KEY (recipe_id)
        REFERENCES public.recipe (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.step
    OWNER to postgres;