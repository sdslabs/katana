package mysql

import (
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sdslabs/katana/lib/utils"
)

func CreateDatabase(database string) error {
	_, err := db.Exec("CREATE DATABASE " + database)
	if err != nil {
		return err
	}
	return nil
}

func CreateGogsUser(username, password string) error {
	// INSERT INTO `user` (`id`, `lower_name`, `name`, `full_name`, `email`, `passwd`, `login_source`, `login_name`, `type`, `location`, `website`, `rands`, `salt`, `created_unix`, `updated_unix`, `last_repo_visibility`, `max_repo_creation`, `is_active`, `is_admin`, `allow_git_hook`, `allow_import_local`, `prohibit_login`, `avatar`, `avatar_email`, `use_custom_avatar`, `num_followers`, `num_following`, `num_stars`, `num_repos`, `description`, `num_teams`, `num_members`) VALUES
	// (NULL, 'hashkat', 'hashkat', '', 'cyberboyrocks@gmail.com', 'eea624a36bf5a06fab84daead5f75ca6ed03061f2bc1864ca60d40185cd93d0962358fd965b3c7438aefc61a14a8b738c29c', '0', '', '0', '', '', 'RwbtuNX5VD', 'mTfyIZC2eV', '1677083612', '1677174739', '0', '-1', '1', '0', '0', '0', '0', '899f8cb105c3866d0e674834320cd7a5', 'cyberboyrocks@gmail.com', '0', '0', '0', '0', '2', '', '0', '0')
	// Get time in unix format and convert it to string
	createdTime := strconv.FormatInt(time.Now().Unix(), 10)
	rand, err := utils.RandomSalt()
	if err != nil {
		return err
	}

	salt, err := utils.RandomSalt()
	if err != nil {
		return err
	}

	password = utils.EncodePassword(password, salt)

	_, error := db.Exec("INSERT INTO `user` (`id`, `lower_name`, `name`, `full_name`, `email`, `passwd`, `login_source`, `login_name`, `type`, `location`, `website`, `rands`, `salt`, `created_unix`, `updated_unix`, `last_repo_visibility`, `max_repo_creation`, `is_active`, `is_admin`, `allow_git_hook`, `allow_import_local`, `prohibit_login`, `avatar`, `avatar_email`, `use_custom_avatar`, `num_followers`, `num_following`, `num_stars`, `num_repos`, `description`, `num_teams`, `num_members`) VALUES (NULL, `" + username + "`), '" + username + "', '', '" + username + "@" + "katana.com" + "', '" + password + "', '0', '', '0', '', '', '" + rand + "', '" + salt + "', '" + createdTime + "', '" + createdTime + "', '0', '-1', '1', '0', '0', '0', '0', '" + utils.MD5(username+"@"+"katana.com") + "', '" + username + "@" + "katana.com" + "', '0', '0', '0', '0', '0', '', '0', '0')")
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
