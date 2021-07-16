package seeder

import (
	"ImaginatoGolangTestTask/shared/utils"
	_const "ImaginatoGolangTestTask/shared/utils/const"
)

func (s Seed) AdminSeed() {
	adminData := []map[string]interface{}{
		{
			"Name":           "super admin",
			"Email":          "superadmin@gmail.com",
			"Password":       utils.HashedPassword("123456"),
			"ResetToken":     "",
			"VerifiedStatus": _const.Verified,
			"Status":         _const.Active,
		},
	}

	for i := 0; i < len(adminData); i++ {
		//prepare the statement
		stmt, _ := s.db.Prepare(`INSERT INTO admin(name, email, password, reset_token, verified_status, status) VALUES (?,?,?,?,?,?)`)
		// execute query
		_, err := stmt.Exec(adminData[i]["Name"], adminData[i]["Email"], adminData[i]["Password"], adminData[i]["ResetToken"], adminData[i]["VerifiedStatus"], adminData[i]["Status"])
		if err != nil {
			panic(err)
		}
	}
}
