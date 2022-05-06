DROP DATABASE IF EXISTS instanTex_db;

CREATE DATABASE instanTex_db;

USE instanTex_db;

CREATE TABLE Users (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    date_of_join DATE NOT NULL,
    state VARCHAR(200) NOT NULL,
    salt INT NOT NULL,
    pHash CHAR(64) NOT NULL,
	last_login TIMESTAMP NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE FriendRequests (
    id INT NOT NULL AUTO_INCREMENT,
    senderId INT NOT NULL,
    reciverId INT NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE PrivateMessages (
    id INT NOT NULL,
    user INT NOT NULL
);

CREATE TABLE GroupMembers (
    id INT NOT NULL,
    user INT NOT NULL
);

CREATE TABLE GroupNames (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(40) NOT NULL,
    description VARCHAR(250),
    PRIMARY KEY ( id )
);

CREATE TABLE Messages (
    id INT NOT NULL AUTO_INCREMENT,
    conv TEXT NOT NULL, -- specifies the group or pm where the message was written (P for pm's, G for groups)
    content VARCHAR(300) NOT NULL,
    attachment VARCHAR(150),
    PRIMARY KEY ( id )
);

CREATE TABLE Token (
    userid INT NOT NULL,
    accessToken CHAR(64) NOT NULL,
    act_expt INT NOT NULL,
    refreshToken CHAR(64) NOT NULL,
	rft_expt INT NOT NULL,
    PRIMARY KEY ( userid )
);

-- CREATE TABLE AccessToken (
-- 	userid INT NOT NULL,
-- 	accessToken CHAR(64) NOT NULL,
-- 	act_expt INT NOT NULL,
-- 	PRIMARY KEY ( userid )
-- );

-- CREATE TABLE RefreshToken (
-- 	Id INT AUTO_INCREMENT
-- 	userid INT NOT NULL,
-- 	refreshToken CHAR(64) NOT NULL,
-- 	rft_expt INT NOT NULL,
-- 	PRIMARY KEY ( userid )
-- );


CREATE TABLE PwOtp (
	userId INT NOT NULL,
	otp CHAR(32) NOT NULL,
	expt INT NOT NULL,
	PRIMARY KEY (userId)
);

CREATE TABLE LoginState (
    idstring CHAR(15) NOT NULL,
    userEmail VARCHAR(45),
    PRIMARY KEY ( idstring )
);


INSERT INTO Users ( username, email, date_of_join, state, salt, pHash, last_login ) VALUES ( "pima", "pippo.mario@gimelli.com", CURRENT_DATE(), "io essere pippo", 66858, "ac7c27a867f92dbecc637a14afae8657f2c2a65eb47faeb3a6cadcad21c17da0", CURRENT_TIMESTAMP() );

INSERT INTO Users ( username, email, date_of_join, state, salt, pHash, last_login) VALUES ( "taurone", "taurone.mario@gimelli.com", CURRENT_DATE(), "no job gang", 66858, "ac7c27a867f92dbecc637a14afae8657f2c2a65eb47faeb3a6cadcad21c17da0", CURRENT_TIMESTAMP());

INSERT INTO PrivateMessages ( id, user ) VALUES ( 1, 1 );
INSERT INTO PrivateMessages ( id, user ) VALUES ( 1, 2 );

INSERT INTO GroupNames ( name, description ) VALUES ( "gruppo tennici", "taurone sei un grande" );
INSERT INTO GroupNames ( name, description ) VALUES ( "gruppo tennici (senza taurone)", "tutti i miei amici odiano taurone" );

INSERT INTO GroupMembers ( id, user ) VALUES ( 1, 1 );
-- INSERT INTO GroupMembers ( id, user ) VALUES ( 1, 2 );
INSERT INTO GroupMembers ( id, user ) VALUES ( 2, 1 );