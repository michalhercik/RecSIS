SET search_path TO recommender;

DROP TABLE IF EXISTS povinn CASCADE;

CREATE TABLE povinn (
    povinn     VARCHAR(10),
    pnazev     VARCHAR(250),
    panazev    VARCHAR(250),
    vplatiod   INT,
    vplatido   INT,
    pfakulta   VARCHAR(5),
    pgarant    VARCHAR(10),
    pvyucovan  VARCHAR(1),
    vsemzac    VARCHAR(1),
    vsempoc    INT,
    vrozsahpr1 INT,
    vrozsahcv1 INT,
    vrozsahpr2 INT,
    vrozsahcv2 INT,
    vrvcem     VARCHAR(2),
    vtyp       VARCHAR(2),
    vebody     INT,
    vucit1     VARCHAR(10),
    vucit2     VARCHAR(10),
    vucit3     VARCHAR(10)
);

DROP TABLE IF EXISTS searchable_povinn CASCADE;

CREATE TABLE searchable_povinn (
    povinn     VARCHAR(10)
);

DROP TABLE IF EXISTS preq CASCADE;

CREATE TABLE preq (
	POVINN VARCHAR(10),
	REQTYP VARCHAR(1),
	REQPOVINN VARCHAR(10)
);

DROP TABLE IF EXISTS pamela CASCADE;

CREATE TABLE pamela (
	POVINN VARCHAR(10),
	TYP VARCHAR(1),
	JAZYK VARCHAR(6),
	MEMO TEXT
);

DROP TABLE IF EXISTS studium CASCADE;

CREATE TABLE studium (
    soident INT,
    sident INT,
    sfak VARCHAR(5),
    sfak2 VARCHAR(5),
    sdruh VARCHAR(2),
    sobor VARCHAR(12),
    srokp VARCHAR(4),
    sstav VARCHAR(6),
    sroc INT,
    splan VARCHAR(15)
);

DROP TABLE IF EXISTS zkous CASCADE;

CREATE TABLE zkous (
    zident INT,
    zskr VARCHAR(4),
    zsem VARCHAR(1),
    zpovinn VARCHAR(10),
    zmarx VARCHAR(5),
    zroc INT,
    zbody INT,
    zsplcelk VARCHAR(1)
);

DROP TABLE IF EXISTS stud_plan CASCADE;

CREATE TABLE stud_plan (
    code VARCHAR(10),
    interchangeability VARCHAR(10),
    bloc_subject_code VARCHAR(20),
    bloc_type VARCHAR(1),
    bloc_grade VARCHAR(50),
    bloc_limit INT,
    bloc_name_cz VARCHAR(250),
    bloc_name_en VARCHAR(250),
    plan_code VARCHAR(15),
    plan_year INT
);

DROP TABLE IF EXISTS obor CASCADE;

CREATE TABLE obor (
    kod VARCHAR(12),
    nazev VARCHAR(250),
    anazev VARCHAR(250)
);

DROP TABLE IF EXISTS trida CASCADE;

CREATE TABLE trida (
    povinn VARCHAR(10),
    kod VARCHAR(7),
    nazev VARCHAR(50)
);

DROP TABLE IF EXISTS klas CASCADE;

CREATE TABLE klas (
    povinn VARCHAR(10),
    kod VARCHAR(6),
    nazev VARCHAR(60),
    anazev VARCHAR(60)
);
