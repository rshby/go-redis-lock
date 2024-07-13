package cache

import "fmt"

func GetStudentCacheKeyByID(id int) string {
	return fmt.Sprintf("object:student:id:%d", id)
}

func GetStudentCacheKeyByEmail(email string) string {
	return fmt.Sprintf("object:student:email:%s", email)
}
