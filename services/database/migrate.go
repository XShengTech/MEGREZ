package database

import "megrez/models"

func autoMigrate() (err error) {
	l.SetFunction("autoMigration")

	err = DB.AutoMigrate(&models.Instances{})
	if err != nil {
		l.Fatal("Failed to migrate containers", err)
	}

	err = DB.AutoMigrate(&models.Orders{})
	if err != nil {
		l.Fatal("Failed to migrate orders", err)
	}

	err = DB.AutoMigrate(&models.Servers{})
	if err != nil {
		l.Fatal("Failed to migrate servers", err)
	}

	err = DB.AutoMigrate(&models.System{})
	if err != nil {
		l.Fatal("Failed to migrate system", err)
	}

	err = DB.AutoMigrate(&models.Users{})
	if err != nil {
		l.Fatal("Failed to migrate users", err)
	}

	return nil
}
