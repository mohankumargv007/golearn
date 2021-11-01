CREATE TABLE employees (
    id int NOT NULL AUTO_INCREMENT,
    emp_id varchar(100) NOT NULL,
    emp_name varchar(100) NOT NULL,
    role text NOT NULL,
    created datetime,
    updated datetime,
    PRIMARY KEY (id)
)