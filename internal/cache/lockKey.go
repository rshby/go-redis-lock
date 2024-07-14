package cache

import "fmt"

func CreateNewStudentLockKeyByEmail(email string) string {
	return fmt.Sprintf("student:create_new_student:email:%s", email)
}
