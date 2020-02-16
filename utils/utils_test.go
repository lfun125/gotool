package utils

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	fmt.Println(CheckEmail("123@d2.cs"))
}

func TestSendToMail(t *testing.T) {
	err := SendToMail(
		"system@123.com",
		"abc123",
		"system",
		"smtp.gmail.com:465",
		"abc123@outlook.com",
		"test",
		"哈哈",
		false,
	)
	fmt.Println(err)
}
