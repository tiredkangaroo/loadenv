package loadenv

import (
	"os"
	"testing"
)

type ExampleStruct1 struct {
	Username                         string
	Password                         string
	TwoFactorAuthenticationSecretKey string
	HelloWorld                       string `required:"false"`
}

func TestLoad(t *testing.T) {
	err := Load(".env2")
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	expected1 := "https://user1:password123@example.com:5432/mydatabase"
	t.Logf("got: %v", os.Getenv("DatabaseConnectionURI"))
	t.Logf("expected: %v", expected1)
	if os.Getenv("DatabaseConnectionURI") != expected1 {
		t.Errorf("got and expected1 did not match.")
		t.FailNow()
	}
}

func TestUnmarshal(t *testing.T) {
	var es1 ExampleStruct1
	err := Unmarshal(&es1, ".env1")
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	expected := ExampleStruct1{
		Username:                         "user1",
		Password:                         "password123",
		TwoFactorAuthenticationSecretKey: "dab684d4a86cda8cdb4cdbac6b8adcb4a68dcb",
		HelloWorld:                       "",
	}
	t.Logf("got: %v", es1)
	t.Logf("expected: %v", expected)
	if es1 != expected {
		t.Errorf("got and expected did not match.")
		t.FailNow()
	}
}
