package cache

import "fmt"

func CreateNewStudentLockKeyByEmail(email string) string {
	return fmt.Sprintf("student:create_new_student:email:%s", email)
}

func BurstStudentCountLockKeyByEpoch(epoch int64) string {
	return fmt.Sprintf("student:burst_student_count:epoch:%d", epoch)
}
