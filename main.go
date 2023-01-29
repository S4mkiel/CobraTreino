package main

// Importing the packages
import (
	f "fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/cobra"
)

// Creating the User
type User struct {
	gorm.Model
	Username     string `gorm:"unique_index;not null"`
	Name         string
	Age          uint
	CompanyID    uint    `gorm:"ForeignKey:CompanyRefer"`
	CompanyRefer Company `gorm:"ForeignKey:CompanyID;AssociationForeignKey:ID"`
}

// Creating the struct Company
type Company struct {
	gorm.Model
	Name string `gorm:"unique_index;not null"`
}

// Main function
func main() {
	// Creating the database
	db, err := gorm.Open("sqlite3", "test.db")
	// Error handling
	if err != nil {
		f.Println(err)
		return
	}
	// Closing the database
	defer db.Close()
	// AutoMigrate the database
	if err := db.AutoMigrate(&User{}, &Company{}).Error; err != nil {
		log.Fatal(err)
		return
	}
	// Creating the commands using cobra
	var rootCmd = &cobra.Command{Use: "app"}
	// Adding the commands to the root command
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a user and company",
		Run: func(cmd *cobra.Command, args []string) {
			Username, _ := cmd.Flags().GetString("username")
			Name, _ := cmd.Flags().GetString("name")
			Age, _ := cmd.Flags().GetUint("age")
			nameCompany, _ := cmd.Flags().GetString("nameCompany")
			companyID, _ := cmd.Flags().GetUint("companyID")
			// Creating the transaction to the database
			tx := db.Begin()
			if err := tx.Create(&Company{Name: nameCompany}).Error; err != nil { // Creating the company
				f.Println("Error creating company:", err)
				tx.Rollback() // Rollback the transaction
				return        // Return the error
			}
			if err := tx.Create(&User{Username: Username, Name: Name, Age: Age, CompanyID: companyID}).Error; err != nil { // Creating the user
				f.Println("Error creating user:", err)
				tx.Rollback() // Rollback the transaction
				return        // Return the error
			}
			tx.Commit()
			f.Println("User and Company created with successfully.") // Print the message if the user and company was created with successfully
		},
	}
	// Adding the flags to the update command
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a user and company",
		Run: func(cmd *cobra.Command, args []string) { // Run the command
			Username, _ := cmd.Flags().GetString("username")
			Name, _ := cmd.Flags().GetString("name")
			Age, _ := cmd.Flags().GetUint("age")
			nameCompany, _ := cmd.Flags().GetString("nameCompany")
			CompanyID, _ := cmd.Flags().GetUint("companyID")
			tx := db.Begin()
			var user User
			// Searching the user and company
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
	var deleteCmd = &cobra.Command{
		Use:   "delete",
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
	var listCmd = &cobra.Command{
		Use:   "search",
		Short: "search a user and company",
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

	deleteCmd.Flags().String("username", "", "delete username for user")
	deleteCmd.MarkFlagRequired("username")
	deleteCmd.Flags().String("nameCompany", "", "delete name for company")
	deleteCmd.MarkFlagRequired("nameCompany")

	listCmd.Flags().String("username", "", "search username for user")
	listCmd.MarkFlagRequired("username")
	listCmd.Flags().String("nameCompany", "", "search name for company")
	listCmd.MarkFlagRequired("nameCompany")
	rootCmd.AddCommand(createCmd, updateCmd, deleteCmd, listCmd)
	rootCmd.Execute()
}
