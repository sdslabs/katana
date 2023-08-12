package mysql

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
)

func CreateDatabase(database string) error {
	_, err := db.Exec("CREATE DATABASE " + database)
	if err != nil {
		return err
	}
	return nil
}

func CreateGogsUser(username, password string) error {

	// Get time in unix format and convert it to string
	gogs, err := sql.Open("mysql", configs.MySQLConfig.Username+":"+configs.MySQLConfig.Password+"@tcp("+utils.GetKatanaLoadbalancer()+":3306)/gogs")
	if err != nil {
		return err
	}
	createdTime := time.Now().Unix()
	rand, err := utils.RandomSalt()
	if err != nil {
		log.Println(err)
	}

	salt, err := utils.RandomSalt()
	if err != nil {
		log.Println(err)
	}

	password = utils.EncodePassword(password, salt)

	// TODO: use an illegal TLD instead of .com
	user := &types.GogsUser{
		LowerName:   username,
		Name:        username,
		FullName:    username,
		Email:       username + "@" + "katana.com",
		Password:    password,
		Rands:       rand,
		Salt:        salt,
		CreatedUnix: createdTime,
		UpdatedUnix: createdTime,
		Avatar:      utils.MD5(username + "@" + "katana.com"),
		AvatarEmail: username + "@" + "katana.com",
	}

	query := "INSERT INTO `user` (`id`, `lower_name`, `name`, `full_name`, `email`, `passwd`, `login_source`, `login_name`, `type`, `location`, `website`, `rands`, `salt`, `created_unix`, `updated_unix`, `last_repo_visibility`, `max_repo_creation`, `is_active`, `is_admin`, `allow_git_hook`, `allow_import_local`, `prohibit_login`, `avatar`, `avatar_email`, `use_custom_avatar`, `num_followers`, `num_following`, `num_stars`, `num_repos`, `description`, `num_teams`, `num_members`) VALUES (NULL, '" + user.LowerName + "', '" + user.Name + "', '" + user.FullName + "', '" + user.Email + "', '" + user.Password + "', '0', '', '0', '', '', '" + user.Rands + "', '" + user.Salt + "', '" + strconv.FormatInt(user.CreatedUnix, 10) + "', '" + strconv.FormatInt(user.UpdatedUnix, 10) + "', '0', '-1', '1', '0', '0', '0', '0', '" + user.Avatar + "', '" + user.AvatarEmail + "', '0', '0', '0', '0', '0', '', '0', '0')"
	_, err = gogs.Exec(query)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func CreateGogsAdmin(username, password string) error {
	// Get time in unix format and convert it to string
	gogs, err := sql.Open("mysql", configs.MySQLConfig.Username+":"+configs.MySQLConfig.Password+"@tcp("+utils.GetKatanaLoadbalancer()+":3306)/gogs")
	if err != nil {
		return err
	}

	// Check if an admin already exists
	var adminCount int
	err = gogs.QueryRow("SELECT COUNT(*) FROM `user` WHERE is_admin = 1").Scan(&adminCount)
	if err != nil {
		return err
	}

	if adminCount > 0 {
		return nil
	}

	createdTime := time.Now().Unix()
	rand, err := utils.RandomSalt()
	if err != nil {
		log.Println(err)
	}

	salt, err := utils.RandomSalt()
	if err != nil {
		log.Println(err)
	}

	password = utils.EncodePassword(password, salt)

	user := &types.GogsUser{
		LowerName:   username,
		Name:        username,
		FullName:    username,
		Email:       username + "@" + "katana.com",
		Password:    password,
		Rands:       rand,
		Salt:        salt,
		CreatedUnix: createdTime,
		UpdatedUnix: createdTime,
		Avatar:      utils.MD5(username + "@" + "katana.com"),
		AvatarEmail: username + "@" + "katana.com",
		IsAdmin:     true,
	}

	_, err = gogs.Exec("INSERT INTO `user` (`id`, `lower_name`, `name`, `full_name`, `email`, `passwd`, `login_source`, `login_name`, `type`, `location`, `website`, `rands`, `salt`, `created_unix`, `updated_unix`, `last_repo_visibility`, `max_repo_creation`, `is_active`, `is_admin`, `allow_git_hook`, `allow_import_local`, `prohibit_login`, `avatar`, `avatar_email`, `use_custom_avatar`, `num_followers`, `num_following`, `num_stars`, `num_repos`, `description`, `num_teams`, `num_members`) VALUES (NULL, '" + user.LowerName + "', '" + user.Name + "', '" + user.FullName + "', '" + user.Email + "', '" + user.Password + "', '0', '', '0', '', '', '" + user.Rands + "', '" + user.Salt + "', '" + strconv.FormatInt(user.CreatedUnix, 10) + "', '" + strconv.FormatInt(user.UpdatedUnix, 10) + "', '0', '-1', '1', '1', '0', '0', '0', '" + user.Avatar + "', '" + user.AvatarEmail + "', '0', '0', '0', '0', '0', '', '0', '0')")
	if err != nil {
		log.Println(err)
	}

	return nil
}

func CreateAccessToken(username, token string) error {
	// TODO: Don't create connections again
	gogs, err := sql.Open("mysql", configs.MySQLConfig.Username+":"+configs.MySQLConfig.Password+"@tcp("+utils.GetKatanaLoadbalancer()+":3306)/gogs")
	if err != nil {
		return err
	}

	sha256 := utils.SHA256(token)
	// First 40 characters of sha256
	sha1 := string(sha256[:40])

	var uid int
	// TODO: Don't concatenate strings to form query
	err = gogs.QueryRow("SELECT id FROM user WHERE name = '" + username + "'").Scan(&uid)
	if err != nil {
		return err
	}

	// Get time in unix format and convert it to string
	createdTime := time.Now().Unix()

	accessToken := &types.GogsAccessToken{
		Name:        username,
		Sha1:        sha1,
		Sha256:      sha256,
		UserId:      uid,
		CreatedUnix: createdTime,
		UpdatedUnix: createdTime,
	}

	_, err = gogs.Exec("INSERT INTO access_token (`id`, `name`, `sha1`, `sha256`, `uid`, `created_unix`, `updated_unix`) VALUES (NULL, '" + accessToken.Name + "', '" + accessToken.Sha1 + "', '" + accessToken.Sha256 + "', '" + strconv.Itoa(accessToken.UserId) + "', '" + strconv.FormatInt(accessToken.CreatedUnix, 10) + "', '" + strconv.FormatInt(accessToken.UpdatedUnix, 10) + "')")
	if err != nil {
		return err
	}
	return nil
}
