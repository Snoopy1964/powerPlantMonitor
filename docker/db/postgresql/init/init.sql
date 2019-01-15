CREATE TABLE sensor (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    serial_no character varying(50) NOT NULL,
    unit_type character varying(50) NOT NULL,
    max_safe_value double precision NOT NULL,
    min_safe_value double precision NOT NULL
);


-- ALTER TABLE sensor OWNER TO postgres;

--
-- TOC entry 174 (class 1259 OID 40976)
-- Name: sensor_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE sensor_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


-- ALTER TABLE sensor_id_seq OWNER TO postgres;

--
-- TOC entry 2015 (class 0 OID 0)
-- Dependencies: 174
-- Name: sensor_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

-- ALTER SEQUENCE sensor_id_seq OWNED BY sensor.id;


--
-- TOC entry 173 (class 1259 OID 40970)
-- Name: sensor_reading; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sensor_reading (
    id integer NOT NULL,
    value double precision NOT NULL,
    sensor_id integer,
    taken_on timestamp with time zone
);


--ALTER TABLE sensor_reading OWNER TO postgres;

--
-- TOC entry 172 (class 1259 OID 40968)
-- Name: sensor_reading_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE sensor_reading_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--ALTER TABLE sensor_reading_id_seq OWNER TO postgres;

--
-- TOC entry 2018 (class 0 OID 0)
-- Dependencies: 172
-- Name: sensor_reading_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

-- ALTER SEQUENCE sensor_reading_id_seq OWNED BY sensor_reading.id;


--
-- TOC entry 1888 (class 2604 OID 40981)
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY sensor ALTER COLUMN id SET DEFAULT nextval('sensor_id_seq'::regclass);


--
-- TOC entry 1887 (class 2604 OID 40973)
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY sensor_reading ALTER COLUMN id SET DEFAULT nextval('sensor_reading_id_seq'::regclass);


--
-- TOC entry 2005 (class 0 OID 40978)
-- Dependencies: 175
-- Data for Name: sensor; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sensor (id, name, serial_no, unit_type, max_safe_value, min_safe_value) FROM stdin;
1	boiler_pressure_out	MPR-728	MPa	15.4	15.1
4	condensor_pressure_out	MPR-317	MPa	0.0022000000000000001	0.00080000000000000004
5	turbine_pressure_out	MPR-492	MPa	1.3999999999999999	0.80000000000000004
6	boiler_temp_out	XTLR-145	C	625	580
7	turbine_temp_out	XTLR-145	C	115	98
8	condensor_temp_out	XTLR-145	C	98	83
\.


--
-- TOC entry 2020 (class 0 OID 0)
-- Dependencies: 174
-- Name: sensor_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--
