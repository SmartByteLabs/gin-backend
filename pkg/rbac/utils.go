package rbac

import (
	"strings"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func getLevelFromAccessMap(accessMap []AccessMap) []string {
	set := utils.NewSet[string]()

	for _, am := range accessMap {
		slice := strings.Split(am.RequiredData["level"], ",")
		for _, level := range slice {
			set.Add(level)
		}
	}

	if set.Size() == 0 {
		return database.AllFields
	}

	return set.ToSlice()
}

func encryptPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func comparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
