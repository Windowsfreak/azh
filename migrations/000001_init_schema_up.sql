CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(100),
    training_type VARCHAR(100),
    weekday VARCHAR(20) NOT NULL,
    start_time VARCHAR(10),
    end_time VARCHAR(10),
    first_schedule DATE,
    last_schedule DATE,
    trainer_names TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE members (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    sign_up_date DATE,
    cancellation_date DATE,
    age INTEGER,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE member_courses (
    id SERIAL PRIMARY KEY,
    member_id INTEGER NOT NULL,
    course_id VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE participations (
    id SERIAL PRIMARY KEY,
    member_id INTEGER NOT NULL,
    course_id VARCHAR(10) NOT NULL,
    date VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_member_courses_member_id ON member_courses(member_id);
CREATE INDEX idx_member_courses_course_id ON member_courses(course_id);
CREATE INDEX idx_participations_member_id ON participations(member_id);
CREATE INDEX idx_participations_course_id ON participations(course_id);
CREATE INDEX idx_participations_date ON participations(date);
CREATE INDEX idx_participations_course_date ON participations(course_id, date); -- For queries filtering by course_id and date
CREATE INDEX idx_participations_date ON participations(date); -- For range queries on date (export)