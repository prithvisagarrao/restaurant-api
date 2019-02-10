--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.6
-- Dumped by pg_dump version 9.5.6

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: ratings; Type: TABLE; Schema: public; Owner: hellofresh
--

CREATE TABLE ratings (
    rating_id integer NOT NULL,
    recipe_id integer,
    rate integer
);


ALTER TABLE ratings OWNER TO hellofresh;

--
-- Name: ratings_rating_id_seq; Type: SEQUENCE; Schema: public; Owner: hellofresh
--

CREATE SEQUENCE ratings_rating_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE ratings_rating_id_seq OWNER TO hellofresh;

--
-- Name: ratings_rating_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: hellofresh
--

ALTER SEQUENCE ratings_rating_id_seq OWNED BY ratings.rating_id;


--
-- Name: recipes; Type: TABLE; Schema: public; Owner: hellofresh
--

CREATE TABLE recipes (
    recipe_id integer NOT NULL,
    name text NOT NULL,
    prep_time text,
    difficulty integer DEFAULT 1,
    vegetarian boolean DEFAULT false
);


ALTER TABLE recipes OWNER TO hellofresh;

--
-- Name: recipes_recipe_id_seq; Type: SEQUENCE; Schema: public; Owner: hellofresh
--

CREATE SEQUENCE recipes_recipe_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE recipes_recipe_id_seq OWNER TO hellofresh;

--
-- Name: recipes_recipe_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: hellofresh
--

ALTER SEQUENCE recipes_recipe_id_seq OWNED BY recipes.recipe_id;


--
-- Name: user_info; Type: TABLE; Schema: public; Owner: hellofresh
--

CREATE TABLE user_info (
    user_id integer NOT NULL,
    email text NOT NULL,
    password text NOT NULL
);


ALTER TABLE user_info OWNER TO hellofresh;

--
-- Name: user_info_user_id_seq; Type: SEQUENCE; Schema: public; Owner: hellofresh
--

CREATE SEQUENCE user_info_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE user_info_user_id_seq OWNER TO hellofresh;

--
-- Name: user_info_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: hellofresh
--

ALTER SEQUENCE user_info_user_id_seq OWNED BY user_info.user_id;


--
-- Name: rating_id; Type: DEFAULT; Schema: public; Owner: hellofresh
--

ALTER TABLE ONLY ratings ALTER COLUMN rating_id SET DEFAULT nextval('ratings_rating_id_seq'::regclass);


--
-- Name: recipe_id; Type: DEFAULT; Schema: public; Owner: hellofresh
--

ALTER TABLE ONLY recipes ALTER COLUMN recipe_id SET DEFAULT nextval('recipes_recipe_id_seq'::regclass);


--
-- Name: user_id; Type: DEFAULT; Schema: public; Owner: hellofresh
--

ALTER TABLE ONLY user_info ALTER COLUMN user_id SET DEFAULT nextval('user_info_user_id_seq'::regclass);


--
-- Name: recipe_id_key; Type: CONSTRAINT; Schema: public; Owner: hellofresh
--

ALTER TABLE ONLY recipes
    ADD CONSTRAINT recipe_id_key PRIMARY KEY (recipe_id);


--
-- Name: recipe_name_unq; Type: CONSTRAINT; Schema: public; Owner: hellofresh
--

ALTER TABLE ONLY recipes
    ADD CONSTRAINT recipe_name_unq UNIQUE (name);


--
-- Name: ratings_recipe_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: hellofresh
--

ALTER TABLE ONLY ratings
    ADD CONSTRAINT ratings_recipe_id_fkey FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--
