package database

import (
	"backend-go/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	if err := seedUser(db); err != nil {
		return err
	}
	if err := seedCategories(db); err != nil {
		return err
	}
	if err := seedContents(db); err != nil {
		return err
	}
	return nil
}

func seedUser(db *gorm.DB) error {
	const email = "admin@admin.com"

	var count int64
	if err := db.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("[seed] usuário admin já existe, pulando seed")
		return nil
	}

	password := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     "Administrador",
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	log.Printf("[seed] usuário admin criado com sucesso (ID=%d)\n", user.ID)
	return nil
}

func seedCategories(db *gorm.DB) error {
	categories := []models.Category{
		{Name: "Drama"},
		{Name: "Ficção Científica"},
		{Name: "Documentário"},
		{Name: "Ação"},
		{Name: "Comédia"},
	}
	for _, category := range categories {
		var existing models.Category
		if err := db.Where("name = ?", category.Name).
			First(&existing).Error; err == nil {
			continue
		}
		if err := db.Create(&category).Error; err != nil {
			return err
		}
		log.Printf("[seed] categoria criada: %s\n", category.Name)
	}
	return nil
}

func seedContents(db *gorm.DB) error {
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	categoryByName := make(map[string]models.Category, len(categories))
	for _, cat := range categories {
		categoryByName[cat.Name] = cat
	}

	entries := []struct {
		Content    models.Content
		Categories []string
	}{
		{
			Content: models.Content{
				Title:       "Horizonte Vermelho",
				Description: "Uma tripulação enfrenta dilemas éticos em missão a Marte.",
				Type:        "Filme",
				ReleaseDate: time.Date(2022, 5, 12, 0, 0, 0, 0, time.UTC),
			},
			Categories: []string{"Ficção Científica", "Drama"},
		},
		{
			Content: models.Content{
				Title:       "Risos em Família",
				Description: "Série leve sobre uma família tentando equilibrar carreira e humor.",
				Type:        "Série",
				ReleaseDate: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
			},
			Categories: []string{"Comédia"},
		},
		{
			Content: models.Content{
				Title:       "Raízes do Agro",
				Description: "Documentário sobre inovação sustentável no campo brasileiro.",
				Type:        "Documentário",
				ReleaseDate: time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			Categories: []string{"Documentário", "Sustentabilidade"},
		},
	}

	for _, entry := range entries {
		var existing models.Content
		if err := db.Where("title = ?", entry.Content.Title).First(&existing).Error; err == nil {
			continue
		}

		content := entry.Content
		if err := db.Create(&content).Error; err != nil {
			return err
		}

		if len(entry.Categories) > 0 {
			var cats []models.Category
			for _, name := range entry.Categories {
				if cat, ok := categoryByName[name]; ok {
					cats = append(cats, cat)
				}
			}
			if len(cats) > 0 {
				if err := db.Model(&content).Association("Categories").Replace(&cats); err != nil {
					return err
				}
			}
		}

		log.Printf("[seed] conteúdo criado: %s\n", content.Title)
	}

	return nil
}
