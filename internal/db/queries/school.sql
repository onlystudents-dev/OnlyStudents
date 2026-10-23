-- name: GetAccountSubjects :many
WITH acc AS (
    SELECT a.student_id, a.teacher_id, a.guardian_id
    FROM accounts a
    WHERE a.id = sqlc.arg(account_id)
),
     student_ids AS (
         SELECT acc.student_id AS id
         FROM acc
         WHERE acc.student_id IS NOT NULL

         UNION

         SELECT ga.student_id AS id
         FROM guardians_access ga
                  JOIN acc ON ga.guardian_id = acc.guardian_id
     ),
     sched AS (
         SELECT bs.subject_id, bs.custom_subject_id, bs.teacher_id, bs.group_id
         FROM base_schedule bs
                  JOIN acc ON bs.teacher_id = acc.teacher_id
         WHERE bs.school_id = sqlc.arg(school_id)

         UNION

         SELECT bs.subject_id, bs.custom_subject_id, bs.teacher_id, bs.group_id
         FROM base_schedule bs
                  JOIN group_members gm ON gm.group_id = bs.group_id
                  JOIN student_ids si ON gm.student_id = si.id
         WHERE bs.school_id = sqlc.arg(school_id)
     )
SELECT s.id AS id, s.subject_name AS subject_name, FALSE AS custom_subject,
       t.id AS teacher_id,
       COALESCE(NULLIF(t.first_name, ''), t.birth_first_name)::text AS teacher_first_name,
    COALESCE(NULLIF(t.last_name, ''), t.birth_last_name)::text AS teacher_last_name,
    g.id AS group_id, g.group_name AS group_name
FROM sched
         JOIN subjects s ON s.id = sched.subject_id
         JOIN teachers t ON t.id = sched.teacher_id
         JOIN groups g ON g.id = sched.group_id

UNION

SELECT cs.id AS id, cs.subject_name AS subject_name, TRUE AS custom_subject,
       t.id AS teacher_id,
       COALESCE(NULLIF(t.first_name, ''), t.birth_first_name)::text AS teacher_first_name,
    COALESCE(NULLIF(t.last_name, ''), t.birth_last_name)::text AS teacher_last_name,
    g.id AS group_id, g.group_name AS group_name
FROM sched
         JOIN custom_subjects cs ON cs.id = sched.custom_subject_id
         JOIN teachers t ON t.id = sched.teacher_id
         JOIN groups g ON g.id = sched.group_id
WHERE cs.school_id = sqlc.arg(school_id)

ORDER BY id;

-- name: IsAccountInSchool :one
SELECT EXISTS (
    SELECT 1
    FROM accounts a
             JOIN group_members gm ON gm.student_id = a.student_id
    WHERE a.id = sqlc.arg(account_id) AND gm.school_id = sqlc.arg(school_id)

    UNION ALL

    SELECT 1
    FROM accounts a
             JOIN teacher_school ts ON ts.teacher_id = a.teacher_id
    WHERE a.id = sqlc.arg(account_id) AND ts.school_id = sqlc.arg(school_id)

    UNION ALL

    SELECT 1
    FROM accounts a
             JOIN guardians_access ga ON ga.guardian_id = a.guardian_id
             JOIN group_members gm ON gm.student_id = ga.student_id
    WHERE a.id = sqlc.arg(account_id) AND gm.school_id = sqlc.arg(school_id)
) AS in_school;