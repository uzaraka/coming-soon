package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func saveEmail(email string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	var emailList EmailList

	// Read existing file if it exists
	if _, err := os.Stat(filename); err == nil {
		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		if err := json.Unmarshal(data, &emailList); err != nil {
			return fmt.Errorf("error parsing JSON: %w", err)
		}
	}

	for _, existingEmail := range emailList.Emails {
		if existingEmail.Email == email {
			return nil
		}
	}
	user := User{Email: email, Date: time.Now()}
	emailList.Emails = append(emailList.Emails, user)

	data, err := json.MarshalIndent(emailList, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}
