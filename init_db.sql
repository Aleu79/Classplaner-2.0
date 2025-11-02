-- CLASSPLANNER DATABASE INIT SCRIPT

--  DROP
DROP VIEW IF EXISTS vw_task_comments;
DROP VIEW IF EXISTS vw_user_progress_summary;
DROP VIEW IF EXISTS vw_task_overview;
DROP VIEW IF EXISTS vw_class_summary;
DROP VIEW IF EXISTS vw_user_submissions;
DROP VIEW IF EXISTS vw_user_ranking;
DROP VIEW IF EXISTS vw_dashboard_classes;
DROP VIEW IF EXISTS vw_user_timeline;
DROP VIEW IF EXISTS vw_class_submissions_stats;
DROP VIEW IF EXISTS vw_dashboard_students;
DROP VIEW IF EXISTS vw_dashboard_admin;

DROP TABLE IF EXISTS address CASCADE;
DROP TABLE IF EXISTS user_class CASCADE;
DROP TABLE IF EXISTS submission CASCADE;
DROP TABLE IF EXISTS comment CASCADE;
DROP TABLE IF EXISTS calendar CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS classes CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS role CASCADE;

DROP TYPE IF EXISTS task_status;

-- ENUM TYPES
CREATE TYPE task_status AS ENUM ('pendiente', 'entregada', 'atrasada');

--  CREATE TABLES

-- ROLES
CREATE TABLE role (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(250) NOT NULL UNIQUE,
    description VARCHAR(250),
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP
);

-- USERS
CREATE TABLE "users" (
    id          BIGSERIAL PRIMARY KEY,
    username    VARCHAR(100) NOT NULL UNIQUE,
    role_id     BIGINT REFERENCES role(id) ON DELETE SET NULL,
    first_name  VARCHAR(100) NOT NULL,
    last_name   VARCHAR(100) NOT NULL,
    email       VARCHAR(250) NOT NULL UNIQUE,
    phone       VARCHAR(250),
    password    VARCHAR(500) NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP
);

-- CLASSES
CREATE TABLE classes (
    id       BIGSERIAL PRIMARY KEY,
    name     VARCHAR(100) NOT NULL,
    profesor BIGINT REFERENCES "users"(id) ON DELETE SET NULL,
    curso    VARCHAR(100),
    color    VARCHAR(100),
    token    VARCHAR(100) UNIQUE
);

-- TASKS
CREATE TABLE tasks (
    id          BIGSERIAL PRIMARY KEY,
    class_id    BIGINT REFERENCES classes(id) ON DELETE CASCADE,
    titulo      VARCHAR(200) NOT NULL,
    description VARCHAR(1000),
    creado      TIMESTAMP DEFAULT NOW(),
    limite      TIMESTAMP,
    estado      task_status DEFAULT 'pendiente'
);

-- CALENDAR
CREATE TABLE calendar (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(100) NOT NULL,
    description VARCHAR(400),
    id_task     BIGINT REFERENCES tasks(id) ON DELETE CASCADE,
    created     TIMESTAMP DEFAULT NOW(),
    deliver     TIMESTAMP,
    class_name  VARCHAR(100),
    curso       VARCHAR(10)
);

-- COMMENTS
CREATE TABLE comment (
    id          BIGSERIAL PRIMARY KEY,
    task_id     BIGINT REFERENCES tasks(id) ON DELETE CASCADE,
    text        VARCHAR(500) NOT NULL,
    user_name   VARCHAR(100),
    user_photo  VARCHAR(300),
    created_at  TIMESTAMP DEFAULT NOW()
);

-- SUBMISSIONS
CREATE TABLE submission (
    id            BIGSERIAL PRIMARY KEY,
    id_user       BIGINT REFERENCES "users"(id) ON DELETE CASCADE,
    id_task       BIGINT REFERENCES tasks(id) ON DELETE CASCADE,
    file          VARCHAR(500),
    comment       VARCHAR(500),
    date          TIMESTAMP DEFAULT NOW(),
    calification  FLOAT,
    feedback      VARCHAR(500),
    photo         VARCHAR(500),
    CONSTRAINT unique_submission UNIQUE (id_user, id_task)
);

-- USER_CLASS (N:M entre usuario y clase)
CREATE TABLE user_class (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES "users"(id) ON DELETE CASCADE,
    class_id    BIGINT REFERENCES classes(id) ON DELETE CASCADE,
    role        VARCHAR(50) DEFAULT 'student' CHECK (role IN ('student','teacher','assistant')),
    created_at  TIMESTAMP DEFAULT NOW()
);

-- ADDRESS
CREATE TABLE address (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT REFERENCES "users"(id) ON DELETE CASCADE,
    address1     VARCHAR(250),
    address2     VARCHAR(250),
    post_code    VARCHAR(250),
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW(),
    deleted_at   TIMESTAMP
);

-- INDEXES
CREATE INDEX idx_users_username ON "users"(username);
CREATE INDEX idx_users_email ON "users"(email);
CREATE INDEX idx_users_role_id ON "users"(role_id);
CREATE INDEX idx_classes_profesor ON classes(profesor);
CREATE INDEX idx_classes_curso ON classes(curso);
CREATE INDEX idx_tasks_class_id ON tasks(class_id);
CREATE INDEX idx_submission_user_id ON submission(id_user);
CREATE INDEX idx_submission_task_id ON submission(id_task);
CREATE INDEX idx_userclass_user_id ON user_class(user_id);
CREATE INDEX idx_userclass_class_id ON user_class(class_id);
CREATE INDEX idx_address_user_id ON address(user_id);

-- TRIGGERS (actualiza updated_at)
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Actualiza updated_at automáticamente
CREATE TRIGGER trg_user_updated       BEFORE UPDATE ON "users"       FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_role_updated       BEFORE UPDATE ON role         FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_address_updated    BEFORE UPDATE ON address      FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_tasks_updated      BEFORE UPDATE ON tasks        FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_classes_updated    BEFORE UPDATE ON classes      FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_submission_updated BEFORE UPDATE ON submission   FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_comment_updated    BEFORE UPDATE ON comment      FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_user_class_updated BEFORE UPDATE ON user_class   FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- TRIGGER: Estado automático de tareas
CREATE OR REPLACE FUNCTION update_task_state()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.estado IS NULL THEN
        NEW.estado := 'pendiente';
    END IF;

    IF NEW.limite IS NOT NULL AND NEW.limite < NOW() AND NEW.estado = 'pendiente' THEN
        NEW.estado := 'atrasada';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_tasks_state
BEFORE INSERT OR UPDATE ON tasks
FOR EACH ROW EXECUTE FUNCTION update_task_state();

-- TRIGGER: Validar calificación entre 0 y 10
CREATE OR REPLACE FUNCTION validate_grade()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.calification IS NOT NULL AND (NEW.calification < 0 OR NEW.calification > 10) THEN
        RAISE EXCEPTION 'La calificación debe estar entre 0 y 10';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_grade
BEFORE INSERT OR UPDATE ON submission
FOR EACH ROW EXECUTE FUNCTION validate_grade();

-- VIEWS
-- Todas las referencias a usuario ahora usan "users"

-- 1. Submissions con user + task
CREATE OR REPLACE VIEW vw_user_submissions AS
SELECT 
    s.id AS submission_id,
    u.id AS user_id,
    u.username,
    u.first_name,
    u.last_name,
    t.id AS task_id,
    t.titulo AS task_title,
    s.comment,
    s.calification,
    s.feedback,
    s.date
FROM submission s
JOIN "users" u ON u.id = s.id_user
JOIN tasks t ON t.id = s.id_task;

-- 2. Resumen de cada clase
CREATE OR REPLACE VIEW vw_class_summary AS
SELECT 
    c.id AS class_id,
    c.name AS class_name,
    CONCAT(u.first_name, ' ', u.last_name) AS profesor,
    c.curso,
    COUNT(DISTINCT uc.user_id) AS total_students,
    COUNT(DISTINCT t.id) AS total_tasks
FROM classes c
LEFT JOIN "users" u ON u.id = c.profesor
LEFT JOIN user_class uc ON uc.class_id = c.id
LEFT JOIN tasks t ON t.class_id = c.id
GROUP BY c.id, c.name, u.first_name, u.last_name, c.curso;

-- 3. Tareas con su clase, profesor y estado
CREATE OR REPLACE VIEW vw_task_overview AS
SELECT 
    t.id AS task_id,
    t.titulo,
    t.estado,
    t.limite,
    c.name AS class_name,
    CONCAT(u.first_name, ' ', u.last_name) AS profesor
FROM tasks t
JOIN classes c ON c.id = t.class_id
LEFT JOIN "users" u ON u.id = c.profesor;

-- 4. Progreso de usuarios
CREATE OR REPLACE VIEW vw_user_progress_summary AS
SELECT 
    u.id AS user_id,
    u.username,
    COUNT(s.id) AS total_submissions,
    ROUND(AVG(s.calification)::NUMERIC, 2) AS avg_grade,
    SUM(CASE WHEN s.calification IS NOT NULL THEN 1 ELSE 0 END) AS graded_submissions
FROM "users" u
LEFT JOIN submission s ON s.id_user = u.id
GROUP BY u.id, u.username;

-- 5. Comentarios por tarea
CREATE OR REPLACE VIEW vw_task_comments AS
SELECT 
    t.id AS task_id,
    t.titulo AS task_title,
    c.id AS comment_id,
    c.text AS comment_text,
    c.user_name,
    c.user_photo,
    c.created_at
FROM comment c
JOIN tasks t ON t.id = c.task_id;

-- 6. Ranking de alumnos por promedio
CREATE OR REPLACE VIEW vw_user_ranking AS
SELECT
    u.id AS user_id,
    u.username,
    u.first_name,
    u.last_name,
    COUNT(s.id) AS submissions_count,
    ROUND(AVG(s.calification)::NUMERIC, 2) AS avg_grade
FROM "users" u
LEFT JOIN submission s ON s.id_user = u.id
GROUP BY u.id, u.username, u.first_name, u.last_name
ORDER BY avg_grade DESC NULLS LAST, submissions_count DESC;

-- 7. Dashboard extendido por clase
CREATE OR REPLACE VIEW vw_dashboard_classes AS
SELECT
    c.id AS class_id,
    c.name AS class_name,
    c.profesor AS teacher_id,
    COUNT(DISTINCT uc.user_id) FILTER (WHERE uc.role = 'student') AS total_students,
    COUNT(DISTINCT t.id) AS total_tasks,
    COUNT(DISTINCT s.id) AS total_submissions,
    ROUND(AVG(s.calification) FILTER (WHERE s.calification IS NOT NULL)::NUMERIC, 2) AS avg_grade,
    COUNT(DISTINCT s.id) FILTER (WHERE t.estado = 'pendiente') AS pending_tasks,
    COUNT(DISTINCT s.id) FILTER (WHERE t.estado = 'entregada') AS submitted_tasks,
    COUNT(DISTINCT s.id) FILTER (WHERE t.estado = 'atrasada') AS late_tasks
FROM classes c
LEFT JOIN user_class uc ON uc.class_id = c.id
LEFT JOIN tasks t ON t.class_id = c.id
LEFT JOIN submission s ON s.id_task = t.id
GROUP BY c.id, c.name, c.profesor
ORDER BY c.name;

-- 8. Vista combinada usuario, tareas y entregas
CREATE OR REPLACE VIEW vw_user_timeline AS
SELECT
    u.id AS user_id,
    u.username,
    u.first_name,
    u.last_name,
    t.id AS task_id,
    t.titulo AS task_title,
    t.description AS task_description,
    t.creado AS task_created,
    t.limite AS task_deadline,
    t.estado AS task_state,
    cal.id AS calendar_id,
    cal.title AS calendar_title,
    cal.deliver AS calendar_deadline,
    s.id AS submission_id,
    s.date AS submission_date,
    s.calification AS submission_grade,
    s.feedback AS submission_feedback,
    CASE
        WHEN s.id IS NOT NULL THEN 'entregada'
        WHEN t.limite < NOW() THEN 'atrasada'
        ELSE 'pendiente'
    END AS effective_state
FROM "users" u
JOIN user_class uc ON uc.user_id = u.id
JOIN classes c ON c.id = uc.class_id
JOIN tasks t ON t.class_id = c.id
LEFT JOIN calendar cal ON cal.id_task = t.id
LEFT JOIN submission s ON s.id_task = t.id AND s.id_user = u.id
ORDER BY u.id, t.limite;

-- 9. Estadísticas de entregas por clase
CREATE OR REPLACE VIEW vw_class_submissions_stats AS
SELECT
    c.id AS class_id,
    c.name AS class_name,
    COUNT(DISTINCT t.id) AS total_tasks,
    COUNT(DISTINCT s.id) AS total_submissions,
    COUNT(DISTINCT CASE WHEN s.date IS NOT NULL AND s.date <= t.limite THEN s.id END) AS on_time_submissions,
    COUNT(DISTINCT CASE WHEN s.date IS NOT NULL AND s.date > t.limite THEN s.id END) AS late_submissions,
    COUNT(DISTINCT CASE WHEN s.id IS NULL THEN u.id END) AS pending_submissions,
    ROUND(AVG(s.calification)::NUMERIC, 2) AS average_grade
FROM classes c
JOIN tasks t ON t.class_id = c.id
JOIN user_class uc ON uc.class_id = c.id
JOIN "users" u ON u.id = uc.user_id
LEFT JOIN submission s ON s.id_task = t.id AND s.id_user = u.id
GROUP BY c.id, c.name
ORDER BY c.name;

-- 10. Dashboard para el admin
CREATE OR REPLACE VIEW vw_dashboard_admin AS
SELECT
    (SELECT COUNT(*) FROM "users") AS total_users,
    (SELECT COUNT(*) FROM classes) AS total_classes,
    (SELECT COUNT(*) FROM tasks) AS total_tasks,
    (SELECT COUNT(*) FROM submission) AS total_submissions,
    ROUND((SELECT AVG(calification)::NUMERIC 
           FROM submission 
           WHERE calification IS NOT NULL), 2) AS avg_grade,
    (SELECT COUNT(*) FROM comment) AS total_comments;

-- 11. Dashboard de alumnos
CREATE OR REPLACE VIEW vw_dashboard_students AS
SELECT
    u.id AS user_id,
    u.username,
    u.first_name,
    u.last_name,
    c.id AS class_id,
    c.name AS class_name,
    COUNT(DISTINCT t.id) AS total_tasks,
    COUNT(DISTINCT s.id) AS total_submissions,
    ROUND(AVG(s.calification) FILTER (WHERE s.calification IS NOT NULL)::NUMERIC, 2) AS avg_grade,
    COUNT(DISTINCT s.id) FILTER (WHERE t.estado = 'pendiente') AS pending_tasks,
    COUNT(DISTINCT s.id) FILTER (WHERE t.estado = 'entregada') AS submitted_tasks,
    COUNT(DISTINCT s.id) FILTER (WHERE t.estado = 'atrasada') AS late_tasks,
    ROUND(
        CASE 
            WHEN COUNT(DISTINCT t.id) = 0 THEN 0
            ELSE (COUNT(DISTINCT s.id)::DECIMAL / COUNT(DISTINCT t.id)::DECIMAL) * 100
        END, 2
    ) AS progress_percent,
    MAX(s.date) AS last_submission_date
FROM "users" u
JOIN user_class uc ON uc.user_id = u.id
JOIN classes c ON c.id = uc.class_id
LEFT JOIN tasks t ON t.class_id = c.id
LEFT JOIN submission s ON s.id_task = t.id AND s.id_user = u.id
WHERE uc.role = 'student'
GROUP BY u.id, u.username, u.first_name, u.last_name, c.id, c.name
ORDER BY u.id, c.name;
