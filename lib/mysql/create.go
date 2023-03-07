package mysql

import (
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	createdTime := time.Now().Unix()
	rand, err := utils.RandomSalt()
	if err != nil {
		return err
	}

	salt, err := utils.RandomSalt()
	if err != nil {
		return err
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
	}

	_, error := db.Exec("INSERT INTO `user` (`id`, `lower_name`, `name`, `full_name`, `email`, `passwd`, `login_source`, `login_name`, `type`, `location`, `website`, `rands`, `salt`, `created_unix`, `updated_unix`, `last_repo_visibility`, `max_repo_creation`, `is_active`, `is_admin`, `allow_git_hook`, `allow_import_local`, `prohibit_login`, `avatar`, `avatar_email`, `use_custom_avatar`, `num_followers`, `num_following`, `num_stars`, `num_repos`, `description`, `num_teams`, `num_members`) VALUES (NULL, '" + user.LowerName + "'), '" + user.Name + "', '" + user.FullName + "', '" + user.Email + "', '" + user.Password + "', '0', '', '0', '', '', '" + user.Rands + "', '" + user.Salt + "', '" + strconv.FormatInt(user.CreatedUnix, 10) + "', '" + strconv.FormatInt(user.UpdatedUnix, 10) + "', '0', '-1', '1', '0', '0', '0', '0', '" + user.Avatar + "', '" + user.AvatarEmail + "', '0', '0', '0', '0', '0', '', '0', '0')")
	if error != nil {
		return error
	}

	return nil
}

func CreateGogsWebhook(database, webhook string) error {
	_, err := db.Exec("GRANT ALL PRIVILEGES ON " + database + ".* TO '" + webhook + "'@'localhost'")
	if err != nil {
		return err
	}
	return nil
}
