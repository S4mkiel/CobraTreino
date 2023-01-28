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
		f.Println(err)
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
				f.Println("Error creating company:", err)
				tx.Rollback()
				return
			}
			if err := tx.Create(&User{Username: Username, Name: Name, Age: Age, CompanyID: companyID}).Error; err != nil {
				f.Println("Error creating user:", err)
				tx.Rollback()
				return
			}
			tx.Commit()
			f.Println("User and Company created with successfully.")
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

	var updateCmd = &cobra.Command{
		Use:   "update user and company",
		Short: "Update a user and company",
		Run: func(cmd *cobra.Command, args []string) {
			Username, _ := cmd.Flags().GetString("username")
			Name, _ := cmd.Flags().GetString("name")
			Age, _ := cmd.Flags().GetUint("age")
			nameCompany, _ := cmd.Flags().GetString("nameCompany")
			CompanyID, _ := cmd.Flags().GetUint("companyID")
			tx := db.Begin()
			var user User
			db.Where("username = ?", Username).First(&user)
			user.Name = Name
			user.Age = Age
			user.CompanyID = CompanyID
			if err := tx.Save(&user).Error; err != nil {
				f.Println("Error updating user:", err)
				tx.Rollback()
				return
			}
			var company Company
			db.Where("name = ?", nameCompany).First(&company)
			company.Name = nameCompany
			if err := tx.Save(&company).Error; err != nil {
				f.Println("Error updating company:", err)
				tx.Rollback()
				return
			}
			tx.Commit()
			f.Println("User and Company updated with successfully.")
		},
	}

	updateCmd.Flags().String("username", "", "update username for user")
	updateCmd.MarkFlagRequired("username")
	updateCmd.Flags().String("name", "", "update name for user")
	updateCmd.MarkFlagRequired("name")
	updateCmd.Flags().Uint("age", 0, "update age for user")
	updateCmd.MarkFlagRequired("age")
	updateCmd.Flags().String("nameCompany", "", "update name for company")
	updateCmd.MarkFlagRequired("nameCompany")
	updateCmd.Flags().Uint("companyID", 0, "update companyID for use")
	updateCmd.MarkFlagRequired("companyID")
	rootCmd.AddCommand(updateCmd)
	rootCmd.Execute()

	var deleteCmd = &cobra.Command{
		Use:   "Delete user and company",
		Short: "Delete a user and company",
		Run: func(cmd *cobra.Command, args []string) {
			Username, _ := cmd.Flags().GetString("username")
			nameCompany, _ := cmd.Flags().GetString("nameCompany")
			tx := db.Begin()
			var user User
			db.Where("username = ?", Username).First(&user)
			if err := tx.Delete(&user).Error; err != nil {
				f.Println("Error deleting user:", err)
				tx.Rollback()
				return
			}
			var company Company
			db.Where("name = ?", nameCompany).First(&company)
			if err := tx.Delete(&company).Error; err != nil {
				f.Println("Error deleting company:", err)
				tx.Rollback()
				return
			}
			tx.Commit()
			f.Println("User and Company deleted with successfully.")
		},
	}

	deleteCmd.Flags().String("username", "", "delete username for user")
	deleteCmd.MarkFlagRequired("username")
	deleteCmd.Flags().String("nameCompany", "", "delete name for company")
	deleteCmd.MarkFlagRequired("nameCompany")
	rootCmd.AddCommand(deleteCmd)
	rootCmd.Execute()

	var listCmd = &cobra.Command{
		Use:   "Search user and company",
		Short: "Search a user and company",
		Run: func(cmd *cobra.Command, args []string) {
			Username, _ := cmd.Flags().GetString("username")
			nameCompany, _ := cmd.Flags().GetString("nameCompany")
			var user User
			db.Where("username = ?", Username).First(&user)
			f.Println("User:", user)
			var company Company
			db.Where("name = ?", nameCompany).First(&company)
			f.Println("Company:", company)
		},
	}
	listCmd.Flags().String("username", "", "search username for user")
	listCmd.MarkFlagRequired("username")
	listCmd.Flags().String("nameCompany", "", "search name for company")
	listCmd.MarkFlagRequired("nameCompany")
	rootCmd.AddCommand(listCmd)
	rootCmd.Execute()
}
