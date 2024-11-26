package helpers

import "fmt"

func ConnectionString(username, password, host, port, name, sslMode string) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		username, password, host, port, name, sslMode)
}
