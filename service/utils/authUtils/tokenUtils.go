package authUtils

import (
	"service/utils"
)

type PicbotTokens struct {
	PICBOT_01_TOKEN string
	PICBOT_02_TOKEN string
	PICBOT_03_TOKEN string
}

func GetAllToken() PicbotTokens {
	utils.PrintObj("GetAllToken", "")

	model := PicbotTokens{
		PICBOT_01_TOKEN: utils.GetEnvPanic("PICBOT_01_TOKEN"),
		PICBOT_02_TOKEN: utils.GetEnvPanic("PICBOT_02_TOKEN"),
		PICBOT_03_TOKEN: utils.GetEnvPanic("PICBOT_03_TOKEN"),
	}

	return model
}

type PicbotTokensOld struct {
	PICBOT_01_TOKEN_OLD string
	PICBOT_02_TOKEN_OLD string
	PICBOT_03_TOKEN_OLD string
}

func GetAllTokenOld() PicbotTokensOld {
	utils.PrintObj("PicbotTokensOld", "")

	model := PicbotTokensOld{
		PICBOT_01_TOKEN_OLD: utils.GetEnvPanic("PICBOT_01_TOKEN_OLD"),
		PICBOT_02_TOKEN_OLD: utils.GetEnvPanic("PICBOT_02_TOKEN_OLD"),
		PICBOT_03_TOKEN_OLD: utils.GetEnvPanic("PICBOT_03_TOKEN_OLD"),
	}

	return model
}
