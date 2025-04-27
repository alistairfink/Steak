--
-- PostgreSQL database dump
--

-- Dumped from database version 11.3 (Debian 11.3-1.pgdg90+1)
-- Dumped by pg_dump version 11.3 (Debian 11.3-1.pgdg90+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: recipe; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.recipe (id, name, "time", type) FROM stdin;
0662e690-aab2-11e9-990d-0242ac170102	Instant Pot Chicken Thighs	01:30:00	instantpot
\.


--
-- Data for Name: equipment; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.equipment (id, recipe_id, name, quantity) FROM stdin;
e05543ca-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	InstantPot	1
e05596f4-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Wooden Spoon	1
e055ea8c-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Trivet	1
\.


--
-- Data for Name: ingredient; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ingredient (id, recipe_id, name, quantity) FROM stdin;
e052b524-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Chicken Things	12
e05346d8-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Potatoes	3
e0539084-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Onion	1
e053d918-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Paprika	1
e05421fc-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Garlic Salt	1
e0546d42-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Oregano	1
e054b68a-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Thyme	1
e05502de-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Chicken Gravy Mix	1
\.


--
-- Data for Name: picture; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.picture (id, recipe_id, image_source, sort_order) FROM stdin;
f34152f2-b1a2-11e9-9add-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	https://i.imgur.com/vgngezQ.jpg	1
\.


--
-- Data for Name: step; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.step (id, recipe_id, content, step_number) FROM stdin;
e0563744-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Peel and cut potatoes into preffered size.	1
e0568a14-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Peel and cut onion to preffered size.	2
e056cf6a-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Spice skin side of chickent thighs with paprika, garlic salt, oregano, and thyme.	3
e0572208-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Set InstantPot to saute and insert butter or oil.	4
e05769de-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Saute chicken thighs for 4 minutes each side.	5
e057bb00-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Remove chicken thighs and turn off InstantPot.	6
e057fda4-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Poor 1.5 cups of water into InstantPot.	7
e058487c-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Scrape bits from bottom of pot with wooden spoon.	8
e0589c78-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Insert onion, potatoe, trivet, then chicken thighs in that order.	9
e058e6ce-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Set InstantPot to poultry for 13 minutes.	10
e0592e36-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Once InstantPot is done cooking quick release and remove chicken, potatoes, and onion.	11
e05976ac-b197-11e9-b2dd-0242ac170102	0662e690-aab2-11e9-990d-0242ac170102	Poor in gravy mix and set InstantPot to saute. Cook to desired thickness.	12
\.


--
-- PostgreSQL database dump complete
--

