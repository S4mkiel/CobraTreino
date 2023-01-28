package main

import (
	f "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/cobra"
)

// ..a
type User struct {
	gorm.Model
	Username     string `gorm:"unique_index"`
	Name         string
	Age          uint
	CompanyID    uint    `gorm:"ForeignKey:CompanyRefer"`
	CompanyRefer Company `gorm:"ForeignKey:CompanyID;AssociationForeignKey:ID"`
}

type Company struct {
	gorm.Model
	Name string `gorm:"unique_index"`
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
			companyID, _ := cmd.Flags().GetUint("companyID")
			tx := db.Begin()
			if err := tx.Create(&Company{Name: nameCompany}).Error; err != nil {
				fmt.Println("Error creating company:", err)
				tx.Rollback()
				return
			}
			if err := tx.Create(&User{Username: Username, Name: Name, Age: Age, CompanyID: companyID}).Error; err != nil {
				fmt.Println("Error creating user:", err)
				tx.Rollback()
				return
			}
			tx.Commit()
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
	createCmd.Flags().Uint("companyID", 0, "companyID for use")
	createCmd.MarkFlagRequired("companyID")
	rootCmd.AddCommand(createCmd)
	rootCmd.Execute()
}
