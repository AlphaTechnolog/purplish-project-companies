package core

import (
    "database/sql"
    "encoding/json"
    "io"
    "net/http"

    "github.com/alphatechnolog/purplish-project-companies/database"
    "github.com/gin-gonic/gin"
)

func getCompanies(d *sql.DB, c *gin.Context) error {
    companies, err := database.GetCompanies(d)
    if err != nil {
        return err
    }

    c.JSON(http.StatusOK, gin.H{"companies": companies})

    return nil
}

func getCompany(d *sql.DB, c *gin.Context) error {
    companyID := c.Param("ID")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Specify company ID"})
		return nil
	}

    company, err := database.GetCompany(d, companyID)
    if err != nil {
        return err
    }

    c.JSON(http.StatusOK, gin.H{"company": company})

    return nil
}

func createCompany(d *sql.DB, c *gin.Context) error {
    bodyContents, err := io.ReadAll(c.Request.Body)
    if err != nil {
        return err
    }

    var createPayload database.CreateCompanyPayload
    if err = json.Unmarshal(bodyContents, &createPayload); err != nil {
        return err
    }

    if err = database.CreateCompany(d, createPayload); err != nil {
        return err
    }

    c.JSON(http.StatusCreated, gin.H{"ok": true})

    return nil
}

func removeCompany(d *sql.DB, c *gin.Context) error {
	companyID := c.Param("ID")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return nil
	}

	if err := database.RemoveCompany(d, companyID); err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})

	return nil
}

func CreateCompaniesRoutes(d *sql.DB, r *gin.RouterGroup) {
    r.GET("/", WrapError(WithDB(d, getCompanies)))
    r.GET("/:ID", WrapError(WithDB(d, getCompany)))
    r.POST("/", WrapError(WithDB(d, createCompany)))
    r.DELETE("/:ID", WrapError(WithDB(d, removeCompany)))
}