package v1

import (
	"context"
	"fmt"
	db_queries "onlystudents/internal/db/store"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Login(c fiber.Ctx, pool *pgxpool.Pool) error {

	var data map[string]interface{}

	json.Unmarshal(c.Body(), &data)

	userRaw, ok := data["user"]

	if !ok {
		return c.SendStatus(400)
	}

	userNumber, ok := userRaw.(int32)

	user := pgtype.Int4{
		Int32: userNumber,
		Valid: true,
	}

	//passwordRaw, ok := data["password"]

	//if !ok {
	//	return c.SendStatus(400)
	//}

	//password, ok := passwordRaw.(string)

	roleRaw, ok := data["role"]

	if !ok {
		return c.SendStatus(400)
	}

	role, ok := roleRaw.(string)

	queries := db_queries.New(pool)

	if role == "student" {
		student, ok := queries.GetAccountByStudentID(context.Background(), user)
		if ok != nil {
			fmt.Println(ok)
		}

		fmt.Println(student)
	}

	return nil
}
