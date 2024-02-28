--
-- PostgreSQL database dump
--

-- Dumped from database version 14.10 (Ubuntu 14.10-1.pgdg22.04+1)
-- Dumped by pg_dump version 14.10 (Ubuntu 14.10-1.pgdg22.04+1)

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
-- Name: check_order_details(); Type: FUNCTION; Schema: public; Owner: abs_app
--

CREATE FUNCTION public.check_order_details() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
declare derived_menu_type_id int;
declare menu_option_exists int;
declare menu_option_value_exists int;
begin

select menu_types.id from menu join menu_types on menu.type_id=menu_types.id where menu.id=NEW.menu_id into derived_menu_type_id;

select count(*) into menu_option_exists from menu_options where menu_options.menu_type_id=derived_menu_type_id and option_id=NEW.menu_option_id;

select count(*) from menu_options where menu_options.menu_type_id=derived_menu_type_id and option_id=NEW.menu_option_id into menu_option_value_exists;

if menu_option_exists = 1 and menu_option_value_exists = 1 then
	return new;
end if;

end;
$$;


ALTER FUNCTION public.check_order_details() OWNER TO abs_app;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: menu; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.menu (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    type_id bigint,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.menu OWNER TO abs_app;

--
-- Name: menu_available_options; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.menu_available_options (
    id integer NOT NULL,
    option text
);


ALTER TABLE public.menu_available_options OWNER TO abs_app;

--
-- Name: menu_option_values; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.menu_option_values (
    option_id bigint,
    id integer NOT NULL,
    value text
);


ALTER TABLE public.menu_option_values OWNER TO abs_app;

--
-- Name: menu_options; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.menu_options (
    menu_type_id integer NOT NULL,
    option_id integer NOT NULL
);


ALTER TABLE public.menu_options OWNER TO abs_app;

--
-- Name: menu_options_id_seq; Type: SEQUENCE; Schema: public; Owner: abs_app
--

CREATE SEQUENCE public.menu_options_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.menu_options_id_seq OWNER TO abs_app;

--
-- Name: menu_options_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: abs_app
--

ALTER SEQUENCE public.menu_options_id_seq OWNED BY public.menu_available_options.id;


--
-- Name: menu_types; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.menu_types (
    id integer NOT NULL,
    type text
);


ALTER TABLE public.menu_types OWNER TO abs_app;

--
-- Name: menu_types_id_seq; Type: SEQUENCE; Schema: public; Owner: abs_app
--

CREATE SEQUENCE public.menu_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.menu_types_id_seq OWNER TO abs_app;

--
-- Name: menu_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: abs_app
--

ALTER SEQUENCE public.menu_types_id_seq OWNED BY public.menu_types.id;


--
-- Name: option_values_id_seq; Type: SEQUENCE; Schema: public; Owner: abs_app
--

CREATE SEQUENCE public.option_values_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.option_values_id_seq OWNER TO abs_app;

--
-- Name: option_values_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: abs_app
--

ALTER SEQUENCE public.option_values_id_seq OWNED BY public.menu_option_values.id;


--
-- Name: order_details; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.order_details (
    order_id uuid NOT NULL,
    menu_id uuid NOT NULL,
    menu_option_id bigint,
    menu_option_value_id bigint,
    quantity bigint
);


ALTER TABLE public.order_details OWNER TO abs_app;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.orders (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    is_completed boolean DEFAULT false,
    completed_at timestamp with time zone,
    user_id uuid
);


ALTER TABLE public.orders OWNER TO abs_app;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.roles (
    id bigint NOT NULL,
    name text
);


ALTER TABLE public.roles OWNER TO abs_app;

--
-- Name: user_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: abs_app
--

CREATE SEQUENCE public.user_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_roles_id_seq OWNER TO abs_app;

--
-- Name: user_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: abs_app
--

ALTER SEQUENCE public.user_roles_id_seq OWNED BY public.roles.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    email text NOT NULL,
    password text NOT NULL,
    role_id bigint DEFAULT 1 NOT NULL
);


ALTER TABLE public.users OWNER TO abs_app;

--
-- Name: variant_values; Type: TABLE; Schema: public; Owner: abs_app
--

CREATE TABLE public.variant_values (
    menu_id uuid NOT NULL,
    option_id bigint NOT NULL,
    option_value_id bigint NOT NULL,
    price bigint
);


ALTER TABLE public.variant_values OWNER TO abs_app;

--
-- Name: menu_available_options id; Type: DEFAULT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu_available_options ALTER COLUMN id SET DEFAULT nextval('public.menu_options_id_seq'::regclass);


--
-- Name: menu_option_values id; Type: DEFAULT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu_option_values ALTER COLUMN id SET DEFAULT nextval('public.option_values_id_seq'::regclass);


--
-- Name: menu_types id; Type: DEFAULT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu_types ALTER COLUMN id SET DEFAULT nextval('public.menu_types_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.user_roles_id_seq'::regclass);


--
-- Data for Name: menu; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.menu (id, name, type_id, created_at, updated_at) FROM stdin;
b983ff25-c532-43ea-aee7-fc75eaa7c2bb	Espresso	1	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
2e2d4f4c-7ebe-464a-b6fb-e99ca990be6a	Americano	1	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
87bdbe7e-9da8-4821-ac38-86d00000cc44	Cappuccino	1	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
6bfe987f-d130-4acb-89d3-7780c2a082ec	Coffee Latte	1	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
dc95630c-af53-41ed-bebe-3e13b5dd3087	Single Origin "Aceh Gayo"	1	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
6fd9affd-51fc-43f8-bb75-2b2688edab40	Single Origin "Toraja"	1	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
bcbf5b01-d7a0-479b-b874-89d631c7bb0a	Affogatto	1	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
7265416d-f787-42a3-be2d-ea2ee0fa9213	Air Mineral	6	2024-01-21 19:39:19.812113+07	2024-01-21 19:39:19.812113+07
3feab325-4c5d-47c0-ae85-c85699553f35	Black Tea	2	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
40029049-0bea-41d8-a33d-f11540dcdf9d	Green Tea	2	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
0a521032-3f94-4a00-8a71-e0f0621ec749	Pandan Tea	2	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
126ad3b5-322d-4fa8-8af8-d58e6276d544	Lychee Tea	2	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
62763862-4dbe-4460-a091-a3d749a09e6a	Sakura Tea	2	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
88643c01-9ccf-4d40-bf9f-ad14c2c2f426	Lemon Tea	2	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
cbfeb280-60e0-4f3d-bcb8-2a52141ea672	Blue Lagoon	3	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
73661629-1a2c-4e47-959a-79c63c8a280f	Coffee Soda	3	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
6dd6805a-9395-4f21-a1db-a7060930ef9c	Not So Pina Colada	3	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
00e1ce84-1759-4883-a473-c8f49a463966	The Bunaken	3	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
5e8233e8-5e7c-416c-a72e-cc44c62867e5	Waterland Iced Tea	3	2024-01-19 08:56:35.169389+07	2024-01-19 08:56:35.169389+07
b738edb4-6b40-4747-aa9c-edaa31f5c2d3	Kopi Susu Gula Aren	7	2024-01-23 21:14:56.395324+07	2024-01-23 21:14:56.395324+07
03a816a2-778f-412c-87f0-af845803846e	Kopi Tetes Vietnam	7	2024-01-23 21:14:56.395324+07	2024-01-23 21:14:56.395324+07
02faff50-df8f-431b-92e4-ea570017774c	Kopi Krim "Macchiato"	7	2024-01-23 21:14:56.395324+07	2024-01-23 21:14:56.395324+07
85f606db-547c-4cbe-b9cc-890110071913	Kopi Bakar a'la ABStudio	7	2024-01-23 21:14:56.395324+07	2024-01-23 21:14:56.395324+07
4cc5533e-3a7d-47dd-825a-816937038249	Kelapa Bakar Latte	7	2024-01-23 21:15:33.999738+07	2024-01-23 21:15:33.999738+07
a2bec9a6-f902-4d12-aae9-249fe4feaf6f	Pandan Bakar Latte	7	2024-01-23 21:15:33.999738+07	2024-01-23 21:15:33.999738+07
faa0559a-e432-4ff6-945a-c1cea4eead3a	The Smurf	7	2024-01-23 21:15:33.999738+07	2024-01-23 21:15:33.999738+07
0075c665-70db-4948-986a-ff1580327e34	Blue Pea Honey Dipper Tea	8	2024-01-23 21:23:02.033199+07	2024-01-23 21:23:02.033199+07
19aed52d-262b-4a2b-affa-62cdea32966c	Coco-Pandan Tea	8	2024-01-23 21:23:02.033199+07	2024-01-23 21:23:02.033199+07
34116ce4-a193-4776-ac2c-090ee0b60328	Lemongrass & Ginger Tea	8	2024-01-23 21:23:02.033199+07	2024-01-23 21:23:02.033199+07
1b5cc9d7-8bdb-4ff3-9db7-181a50d5cca0	Macchiato black Tea	8	2024-01-23 21:23:02.033199+07	2024-01-23 21:23:02.033199+07
18c20762-b4b2-4cae-8531-b5352cf10046	Popping Fruits Tea Sangria	8	2024-01-23 21:23:02.033199+07	2024-01-23 21:23:02.033199+07
71ded19c-fe4b-4a22-88f5-d8ace67efb84	ABStudio Dream Tea	8	2024-01-23 21:23:02.033199+07	2024-01-23 21:23:02.033199+07
de9324dc-ea2e-4bb8-bd68-9865c4b56397	Choco Malt	9	2024-01-23 21:24:09.987634+07	2024-01-23 21:24:09.987634+07
5de5d778-2428-4c99-a287-3d5ce0a0312a	Coffee Baileys	9	2024-01-23 21:24:09.987634+07	2024-01-23 21:24:09.987634+07
3c4be2d1-6c62-4783-8065-d9197f088468	Matcha Green Tea	9	2024-01-23 21:24:09.987634+07	2024-01-23 21:24:09.987634+07
d42c420d-1cb9-4060-884d-3836856b72a4	Red Velvet	9	2024-01-23 21:24:09.987634+07	2024-01-23 21:24:09.987634+07
5b565cf8-c566-46fa-9407-3ccfb40dc299	Taro	9	2024-01-23 21:24:09.987634+07	2024-01-23 21:24:09.987634+07
693dac32-6de7-4846-b139-b1cc17bab0bf	Thai Tea	9	2024-01-23 21:24:09.987634+07	2024-01-23 21:24:09.987634+07
bf5aad5a-4d81-4ab2-ae64-d0dad9b77061	ABS Favourite	4	2024-01-23 21:46:57.195732+07	2024-01-23 21:46:57.195732+07
5b104f3e-8d3f-4683-b9d0-2ff34018de14	Meaty	4	2024-01-23 21:47:02.747995+07	2024-01-23 21:47:02.747995+07
144fce08-f16b-4496-a3bc-f1459ed00301	Urban Farm	4	2024-01-23 21:47:08.219859+07	2024-01-23 21:47:08.219859+07
bf57090b-00e2-424f-a33c-463614cded01	Ayam Sambal Matah	4	2024-01-23 21:47:16.297206+07	2024-01-23 21:47:16.297206+07
f398834f-ede3-4ad2-ba2f-161d4c5ec7be	Rendang Sapi	4	2024-01-23 21:47:22.36944+07	2024-01-23 21:47:22.36944+07
fa47b4a6-dc3d-496c-8191-a5d5c3ea3db1	Roasty Apple	4	2024-01-23 21:47:28.060847+07	2024-01-23 21:47:28.060847+07
ab2a528c-9c5b-45d0-ba7f-ce91a97c6b67	Banana Strawberry	4	2024-01-23 21:47:36.472014+07	2024-01-23 21:47:36.472014+07
580f2a3e-267d-419c-bc95-3ea98e630ec0	Seasonal Fruits Nojito	3	2024-01-24 11:26:40.986212+07	2024-01-24 11:26:40.986212+07
62ac2905-6991-4d22-81c4-6b0ec88b7fcd	Italian Coffee Mocktails	6	2024-01-24 11:35:26.917338+07	2024-01-24 11:35:26.917338+07
575d8350-9b4a-4ae1-a0e5-d4a1a6f47144	Italian Tea Mocktails	6	2024-01-24 11:35:34.635873+07	2024-01-24 11:35:34.635873+07
560e1f98-f497-421c-8493-3e52c93a5e2d	Italian Milk Blend	6	2024-01-24 11:35:41.214941+07	2024-01-24 11:35:41.214941+07
fc6797b3-17fd-4061-ae51-14fa8049e109	Italian Soda	6	2024-01-24 11:35:46.99144+07	2024-01-24 11:35:46.99144+07
f6696b35-750b-4a0c-a11f-1868fe6dfaa9	Minuman ngetes doang	6	2024-02-07 09:29:23.382115+07	2024-02-11 18:41:56.931764+07
\.


--
-- Data for Name: menu_available_options; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.menu_available_options (id, option) FROM stdin;
1	temp
2	blend
35	pizza topping
\.


--
-- Data for Name: menu_option_values; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.menu_option_values (option_id, id, value) FROM stdin;
1	2	hot
1	1	iced
2	3	blend
35	36	regular
1	37	plain
\.


--
-- Data for Name: menu_options; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.menu_options (menu_type_id, option_id) FROM stdin;
1	1
2	1
7	1
8	1
9	2
4	35
6	1
\.


--
-- Data for Name: menu_types; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.menu_types (id, type) FROM stdin;
2	tea
3	mocktails
4	pizza
6	lainnya
7	Artisan Coffee
8	Artisan Tea
9	Artisan Milk Blend
1	kopi
\.


--
-- Data for Name: order_details; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.order_details (order_id, menu_id, menu_option_id, menu_option_value_id, quantity) FROM stdin;
686156e2-993b-4518-ab42-e757335fcd75	ab2a528c-9c5b-45d0-ba7f-ce91a97c6b67	35	36	2
26a33223-da03-4a5c-8bbb-f7fd6944abef	ab2a528c-9c5b-45d0-ba7f-ce91a97c6b67	35	36	2
37b40f41-27fd-4dfb-82ee-6d85a4b20b1f	bf5aad5a-4d81-4ab2-ae64-d0dad9b77061	35	36	1
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.orders (id, created_at, is_completed, completed_at, user_id) FROM stdin;
686156e2-993b-4518-ab42-e757335fcd75	2024-02-20 15:35:39.374745+07	f	\N	fad6c002-cbba-48dc-81d8-d56a17f5428c
37b40f41-27fd-4dfb-82ee-6d85a4b20b1f	2024-02-21 10:23:22.980231+07	f	\N	46e11084-0baa-4d2f-bf52-6f6a93a78619
26a33223-da03-4a5c-8bbb-f7fd6944abef	2024-02-21 10:18:19.674673+07	t	2024-02-21 10:26:54.973271+07	46e11084-0baa-4d2f-bf52-6f6a93a78619
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.roles (id, name) FROM stdin;
1	user
2	admin
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.users (id, name, created_at, updated_at, deleted_at, email, password, role_id) FROM stdin;
46e11084-0baa-4d2f-bf52-6f6a93a78619	Muhammad Rava	2024-02-21 10:12:43.369861+07	2024-02-21 10:12:43.369861+07	\N	anonymous	anonymous	1
754acd9f-0810-409d-a794-cf4290b90272	admin	2024-02-22 09:36:13.919435+07	2024-02-22 09:36:13.918643+07	\N	admin@admin.com	$2a$14$H0.fHbq0r7kIVRfkOLu0cutuhPCmkLjXPGNoWkShWx.TVinyL2pHe	2
07dcd1c2-38a4-4274-848d-1084bbc8683b	Jose	2024-02-19 09:43:37.642367+07	2024-02-19 09:43:37.642367+07	\N	jose@jose.com	$2a$14$wlKfgfe8NND8dyPmaN.8Je4C41S8xHA3F/x8w4kndL9VDXWwzM.uS	1
fad6c002-cbba-48dc-81d8-d56a17f5428c	Abiman You	2024-02-19 10:01:05.918176+07	2024-02-19 10:01:05.918176+07	\N	abim@abim.com	$2a$14$DhwQKFBX8ZVoIJqBokT6guiKeJ063uBpekz.ZM1ISncnhe.xm1/Qe	1
\.


--
-- Data for Name: variant_values; Type: TABLE DATA; Schema: public; Owner: abs_app
--

COPY public.variant_values (menu_id, option_id, option_value_id, price) FROM stdin;
b983ff25-c532-43ea-aee7-fc75eaa7c2bb	1	2	6000
b983ff25-c532-43ea-aee7-fc75eaa7c2bb	1	1	10000
2e2d4f4c-7ebe-464a-b6fb-e99ca990be6a	1	2	7000
2e2d4f4c-7ebe-464a-b6fb-e99ca990be6a	1	1	10000
87bdbe7e-9da8-4821-ac38-86d00000cc44	1	2	10000
87bdbe7e-9da8-4821-ac38-86d00000cc44	1	1	12000
6bfe987f-d130-4acb-89d3-7780c2a082ec	1	2	10000
6bfe987f-d130-4acb-89d3-7780c2a082ec	1	1	12000
dc95630c-af53-41ed-bebe-3e13b5dd3087	1	2	10000
dc95630c-af53-41ed-bebe-3e13b5dd3087	1	1	12000
6fd9affd-51fc-43f8-bb75-2b2688edab40	1	2	10000
bcbf5b01-d7a0-479b-b874-89d631c7bb0a	1	1	12000
6fd9affd-51fc-43f8-bb75-2b2688edab40	1	1	12000
3feab325-4c5d-47c0-ae85-c85699553f35	1	2	5000
3feab325-4c5d-47c0-ae85-c85699553f35	1	1	6000
40029049-0bea-41d8-a33d-f11540dcdf9d	1	2	5000
40029049-0bea-41d8-a33d-f11540dcdf9d	1	1	6000
88643c01-9ccf-4d40-bf9f-ad14c2c2f426	1	2	7000
88643c01-9ccf-4d40-bf9f-ad14c2c2f426	1	1	10000
126ad3b5-322d-4fa8-8af8-d58e6276d544	1	2	10000
126ad3b5-322d-4fa8-8af8-d58e6276d544	1	1	12000
0a521032-3f94-4a00-8a71-e0f0621ec749	1	2	10000
0a521032-3f94-4a00-8a71-e0f0621ec749	1	1	12000
62763862-4dbe-4460-a091-a3d749a09e6a	1	1	12000
62763862-4dbe-4460-a091-a3d749a09e6a	1	2	10000
cbfeb280-60e0-4f3d-bcb8-2a52141ea672	1	1	12000
73661629-1a2c-4e47-959a-79c63c8a280f	1	1	12000
6dd6805a-9395-4f21-a1db-a7060930ef9c	1	1	12000
580f2a3e-267d-419c-bc95-3ea98e630ec0	1	1	12000
00e1ce84-1759-4883-a473-c8f49a463966	1	1	12000
5e8233e8-5e7c-416c-a72e-cc44c62867e5	1	1	12000
bf5aad5a-4d81-4ab2-ae64-d0dad9b77061	35	36	30000
5b104f3e-8d3f-4683-b9d0-2ff34018de14	35	36	35000
144fce08-f16b-4496-a3bc-f1459ed00301	35	36	35000
bf57090b-00e2-424f-a33c-463614cded01	35	36	45000
f398834f-ede3-4ad2-ba2f-161d4c5ec7be	35	36	45000
ab2a528c-9c5b-45d0-ba7f-ce91a97c6b67	35	36	45000
fa47b4a6-dc3d-496c-8191-a5d5c3ea3db1	35	36	45000
b738edb4-6b40-4747-aa9c-edaa31f5c2d3	1	2	10000
b738edb4-6b40-4747-aa9c-edaa31f5c2d3	1	1	12000
03a816a2-778f-412c-87f0-af845803846e	1	1	12000
03a816a2-778f-412c-87f0-af845803846e	1	2	10000
02faff50-df8f-431b-92e4-ea570017774c	1	2	10000
85f606db-547c-4cbe-b9cc-890110071913	1	2	10000
02faff50-df8f-431b-92e4-ea570017774c	1	1	12000
4cc5533e-3a7d-47dd-825a-816937038249	1	1	12000
4cc5533e-3a7d-47dd-825a-816937038249	1	2	10000
a2bec9a6-f902-4d12-aae9-249fe4feaf6f	1	2	10000
a2bec9a6-f902-4d12-aae9-249fe4feaf6f	1	1	12000
faa0559a-e432-4ff6-945a-c1cea4eead3a	1	1	12000
faa0559a-e432-4ff6-945a-c1cea4eead3a	1	2	10000
0075c665-70db-4948-986a-ff1580327e34	1	2	10000
0075c665-70db-4948-986a-ff1580327e34	1	1	12000
19aed52d-262b-4a2b-affa-62cdea32966c	1	1	12000
19aed52d-262b-4a2b-affa-62cdea32966c	1	2	10000
34116ce4-a193-4776-ac2c-090ee0b60328	1	2	10000
1b5cc9d7-8bdb-4ff3-9db7-181a50d5cca0	1	2	10000
18c20762-b4b2-4cae-8531-b5352cf10046	1	2	10000
71ded19c-fe4b-4a22-88f5-d8ace67efb84	1	2	10000
71ded19c-fe4b-4a22-88f5-d8ace67efb84	1	1	12000
34116ce4-a193-4776-ac2c-090ee0b60328	1	1	12000
18c20762-b4b2-4cae-8531-b5352cf10046	1	1	12000
de9324dc-ea2e-4bb8-bd68-9865c4b56397	2	3	12000
5de5d778-2428-4c99-a287-3d5ce0a0312a	2	3	12000
3c4be2d1-6c62-4783-8065-d9197f088468	2	3	12000
d42c420d-1cb9-4060-884d-3836856b72a4	2	3	12000
5b565cf8-c566-46fa-9407-3ccfb40dc299	2	3	12000
693dac32-6de7-4846-b139-b1cc17bab0bf	2	3	12000
62ac2905-6991-4d22-81c4-6b0ec88b7fcd	1	37	15000
575d8350-9b4a-4ae1-a0e5-d4a1a6f47144	1	37	15000
560e1f98-f497-421c-8493-3e52c93a5e2d	1	37	15000
fc6797b3-17fd-4061-ae51-14fa8049e109	1	37	15000
7265416d-f787-42a3-be2d-ea2ee0fa9213	1	37	5000
f6696b35-750b-4a0c-a11f-1868fe6dfaa9	1	37	10000
f6696b35-750b-4a0c-a11f-1868fe6dfaa9	1	1	15000
\.


--
-- Name: menu_options_id_seq; Type: SEQUENCE SET; Schema: public; Owner: abs_app
--

SELECT pg_catalog.setval('public.menu_options_id_seq', 35, true);


--
-- Name: menu_types_id_seq; Type: SEQUENCE SET; Schema: public; Owner: abs_app
--

SELECT pg_catalog.setval('public.menu_types_id_seq', 9, true);


--
-- Name: option_values_id_seq; Type: SEQUENCE SET; Schema: public; Owner: abs_app
--

SELECT pg_catalog.setval('public.option_values_id_seq', 37, true);


--
-- Name: user_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: abs_app
--

SELECT pg_catalog.setval('public.user_roles_id_seq', 1, false);


--
-- Name: users idx_users_email; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT idx_users_email UNIQUE (email);


--
-- Name: users member_pkey; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT member_pkey PRIMARY KEY (id);


--
-- Name: menu_available_options menu_options_pkey; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu_available_options
    ADD CONSTRAINT menu_options_pkey PRIMARY KEY (id);


--
-- Name: menu_types menu_types_pkey; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu_types
    ADD CONSTRAINT menu_types_pkey PRIMARY KEY (id);


--
-- Name: menu minuman_id_minuman_key; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu
    ADD CONSTRAINT minuman_id_minuman_key UNIQUE (id);


--
-- Name: menu_option_values option_values_pkey; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu_option_values
    ADD CONSTRAINT option_values_pkey PRIMARY KEY (id);


--
-- Name: order_details order_details_pkey; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.order_details
    ADD CONSTRAINT order_details_pkey PRIMARY KEY (order_id, menu_id);


--
-- Name: orders pesanan_id_pesanan_key; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT pesanan_id_pesanan_key UNIQUE (id);


--
-- Name: menu pk_minuman; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu
    ADD CONSTRAINT pk_minuman PRIMARY KEY (id);


--
-- Name: orders pk_pesanan; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT pk_pesanan PRIMARY KEY (id);


--
-- Name: roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: variant_values variant_values_pkey; Type: CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.variant_values
    ADD CONSTRAINT variant_values_pkey PRIMARY KEY (menu_id, option_id, option_value_id);


--
-- Name: order_details check_order_details; Type: TRIGGER; Schema: public; Owner: abs_app
--

CREATE TRIGGER check_order_details BEFORE INSERT ON public.order_details FOR EACH ROW EXECUTE FUNCTION public.check_order_details();


--
-- Name: menu fk_menu_type; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu
    ADD CONSTRAINT fk_menu_type FOREIGN KEY (type_id) REFERENCES public.menu_types(id);


--
-- Name: variant_values fk_menu_variant_values; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.variant_values
    ADD CONSTRAINT fk_menu_variant_values FOREIGN KEY (menu_id) REFERENCES public.menu(id);


--
-- Name: orders fk_orders_user; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: users fk_users_role; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: variant_values fk_variant_values_option; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.variant_values
    ADD CONSTRAINT fk_variant_values_option FOREIGN KEY (option_id) REFERENCES public.menu_available_options(id);


--
-- Name: variant_values fk_variant_values_option_value; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.variant_values
    ADD CONSTRAINT fk_variant_values_option_value FOREIGN KEY (option_value_id) REFERENCES public.menu_option_values(id);


--
-- Name: menu_options menu_options_menu_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu_options
    ADD CONSTRAINT menu_options_menu_type_id_fkey FOREIGN KEY (menu_type_id) REFERENCES public.menu_types(id);


--
-- Name: menu_options menu_options_option_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu_options
    ADD CONSTRAINT menu_options_option_id_fkey FOREIGN KEY (option_id) REFERENCES public.menu_available_options(id);


--
-- Name: menu_option_values option_values_option_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.menu_option_values
    ADD CONSTRAINT option_values_option_id_fkey FOREIGN KEY (option_id) REFERENCES public.menu_available_options(id);


--
-- Name: order_details order_details_menu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.order_details
    ADD CONSTRAINT order_details_menu_id_fkey FOREIGN KEY (menu_id) REFERENCES public.menu(id);


--
-- Name: order_details order_details_menu_option_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.order_details
    ADD CONSTRAINT order_details_menu_option_id_fkey FOREIGN KEY (menu_option_id) REFERENCES public.menu_available_options(id);


--
-- Name: order_details order_details_menu_option_value_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.order_details
    ADD CONSTRAINT order_details_menu_option_value_id_fkey FOREIGN KEY (menu_option_value_id) REFERENCES public.menu_option_values(id);


--
-- Name: order_details order_details_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.order_details
    ADD CONSTRAINT order_details_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: variant_values variant_values_menu_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.variant_values
    ADD CONSTRAINT variant_values_menu_id_fkey FOREIGN KEY (menu_id) REFERENCES public.menu(id);


--
-- Name: variant_values variant_values_menu_option_value_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.variant_values
    ADD CONSTRAINT variant_values_menu_option_value_id_fkey FOREIGN KEY (option_value_id) REFERENCES public.menu_option_values(id);


--
-- Name: variant_values variant_values_option_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: abs_app
--

ALTER TABLE ONLY public.variant_values
    ADD CONSTRAINT variant_values_option_id_fkey FOREIGN KEY (option_id) REFERENCES public.menu_available_options(id);


--
-- PostgreSQL database dump complete
--

