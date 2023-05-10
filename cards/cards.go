package cards

type Map map[string]string

func GetData() Map {
	dict := Map{
		"Name":   "Apparition",
		"Image":  "Apparition.png",
		"Color":  "Colorless",
		"Rarity": "Special",
		"Type":   "Skill",
		"Cost":   "1",
		"Text":   "[#Ethereal.|\nGain 1 #Intangible. \n#Exhaust.",
		"Traits": "0",
	}
	return dict
}
