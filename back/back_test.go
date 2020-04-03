package main

import (
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"

	"github.com/devansh42/shree/remote"
)

func TestAuthUserExists(t *testing.T) {
	b, cli := initTestEnvironment()
	uid := time.Now().Unix()
	cli.HMSet(hash(sprint("u", uid)), map[string]interface{}{
		"e": "devansh@gmail.com",
		"p": hash("hello1234"),
		"u": "devansh42",
	})
	defer cli.Del(hash(sprint("u", uid))) //Cleanup
	cli.Set(hash(sprint("u:", "devansh42")), sprint(uid), time.Minute)
	defer cli.Del(hash(sprint("u:", "devansh42"))) //Cleanup
	user := new(remote.User)
	user.Username = "devansh42"
	user.Password = []byte("hello1234")
	resp := new(remote.User)
	err := b.Auth(user, resp)
	if err != nil {
		t.Error(err)
	}
	t.Log("User Authenticated\t", resp)
}

func TestAuthUserDoesntExists(t *testing.T) {
	b, cli := initTestEnvironment()
	user := new(remote.User)
	user.Username = "devansh42"
	user.Password = []byte("hello1234")
	user.Email = "devanshguptamrt@gmail.com"
	resp := new(remote.User)
	err := b.Auth(user, resp)
	if err != nil {
		t.Error(err)
	}
	t.Log("User Authenticated\t", resp)
	//Cleanup
	defer cli.Del(hash(sprint("u", resp.Uid)))
	defer cli.Del(hash(sprint("u:", "devansh42")))
}

func TestAuthUserInvalidtCredentials(t *testing.T) {
	b, cli := initTestEnvironment()
	uid := time.Now().Unix()
	cli.HMSet(hash(sprint("u", uid)), map[string]interface{}{
		"e": "devansh@gmail.com",
		"p": hash("xhello1234"),
		"u": "devansh42",
	})
	defer cli.Del(hash(sprint("u", uid))) //Cleanup
	cli.Set(hash(sprint("u:", "devansh42")), sprint(uid), time.Minute)
	defer cli.Del(hash(sprint("u:", "devansh42"))) //Cleanup

	user := new(remote.User)
	user.Username = "devansh42"
	user.Password = []byte("hello1234")
	resp := new(remote.User)
	err := b.Auth(user, resp)
	if err != nil {
		t.Skip(err)
	}
	t.Log("User Authenticated\t", resp)
}

func initTestEnvironment() (*Backend, *redis.Client) {
	os.Setenv(REDIS_SERVER_ADDR, ":6379")
	cli := redis.NewClient(&redis.Options{Network: "tcp", Addr: ":6379"})
	return new(Backend), cli
}

//Other endpoints are ont required to be tested
