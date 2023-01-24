package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/cobra"
)

type User struct {
	gorm.Model
	Username 		string `gorm:"unique_index"`
	Name  			string
	Age 			uint
	Companies 		[]Company
}

type Company struct {
	gorm.Model
	nameCompany 	string `gorm:"unique_index"`
	UserID 			uint `gorm:"ForeignKey:UserRefer"`
    UserRefer 		User `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`

}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Company{})

	var rootCmd = &cobra.Command{Use: "app"}
	var createCmd = &cobra.Command{
		Use:   "create user and company",
		Short: "Create a user and company",
		Run: func(cmd *cobra.Command, args []string) {
			Username, _ := cmd.Flags().GetString("username") 
			Name, _ := cmd.Flags().GetString("name") 
			Age, _ := cmd.Flags().GetUint("age") 
			nameCompany, _ := cmd.Flags().GetString("nameCompany") 
			db.Create(&User{Username: Username, Name: Name, Age: Age})
			db.Create(&Company{nameCompany: nameCompany})
			fmt.Println("User and Company created with successfully.")
		},
	}
	createCmd.Flags().String("username", "", "username for user")
	createCmd.MarkFlagRequired("username")
	createCmd.Flags().String("name", "", "name for user")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().Uint("age", 0, "age for user")
	createCmd.MarkFlagRequired("age")
	createCmd.Flags().String("nameCompany", "", "name for company")
	createCmd.MarkFlagRequired("nameCompany")
	rootCmd.AddCommand(createCmd)
	rootCmd.Execute()
}