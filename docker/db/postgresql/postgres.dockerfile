FROM postgres

ADD ./init/init.sql /entry-point-initdb.d

