package dotfilesmanager

import (
	"github.com/samber/mo"
)

type DotfilesData struct {
	Email         string                         `mapstructure:"email"`
	FirstName     string                         `mapstructure:"first_name"`
	LastName      string                         `mapstructure:"last_name"`
	GpgSigningKey mo.Option[string]              `mapstructure:"gpg_signing_key"`
	WorkEnv       mo.Option[DotfilesWorkEnvData] `mapstructure:"work_env"`
	SystemData    mo.Option[DotfilesSystemData]  `mapstructure:"system_data"`
}

type DotfilesWorkEnvData struct {
	WorkName  string `mapstructure:"work_name"`
	WorkEmail string `mapstructure:"work_email"`
}

type DotfilesSystemData struct {
	Shell           string `mapstructure:"shell"`
	User            string `mapstructure:"user"`
	MultiUserSystem bool   `mapstructure:"multi_user_system"`
	BrewUser        bool   `mapstructure:"brew_user"`
}
