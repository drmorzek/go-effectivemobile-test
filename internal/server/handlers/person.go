package handlers

import (
	"go-test/internal/db"
	"go-test/internal/server/utils"
	"go-test/pkg/framework"
	"net/http"
)

func (h *Handlers) GetPeople(ctx *framework.Context) {

	name, err := ctx.GetQuery("name")

	if err != nil {
		utils.CheckErrorCtx(ctx, err)
		return
	}

	person := []db.Person{}

	m := make(map[string]interface{})
	m["name"] = name

	h.DB.Where(m).Find(&person)
	ctx.JSON(http.StatusOK, framework.H{
		"result": person,
	})
}

func (h *Handlers) PostPeople(ctx *framework.Context) {

	person := db.Person{
		Name:       ctx.Body["name"].(string),
		Surname:    ctx.Body["surname"].(string),
		Patronymic: ctx.Body["patronymic"].(string),
	}

	find_persons := []db.Person{}
	m := make(map[string]interface{})
	m["name"] = person.Name
	m["surname"] = person.Surname

	h.DB.Where(m).Find(&find_persons)

	if len(find_persons) > 0 {
		ctx.JSON(http.StatusBadRequest, framework.H{
			"error": "Person with this name and surname already yet",
		})
		return
	}

	if err := utils.EnrichPersonData(&person); err != nil {
		utils.CheckErrorCtx(ctx, err)
		return
	}

	if err := h.DB.Omit("ID").Create(&person).Error; err != nil {
		utils.CheckErrorCtx(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, framework.H{
		"result": person,
	})
}

func (h *Handlers) PutPeople(ctx *framework.Context) {

	find_person := db.Person{}

	m := make(map[string]interface{})
	m["id"] = ctx.Params["id"]

	if err := h.DB.Where(m).First(&find_person).Error; err != nil {
		utils.CheckErrorCtx(ctx, err)
		return
	}

	find_person.Name = ctx.Body["name"].(string)
	find_person.Surname = ctx.Body["surname"].(string)
	find_person.Patronymic = ctx.Body["patronymic"].(string)

	if err := h.DB.Save(&find_person).Error; err != nil {
		utils.CheckErrorCtx(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, framework.H{
		"result": find_person,
	})
}
func (h *Handlers) DeletePeople(ctx *framework.Context) {

	find_person := db.Person{}

	m := make(map[string]interface{})
	m["id"] = ctx.Params["id"]

	if err := h.DB.Where(m).First(&find_person).Error; err != nil {
		utils.CheckErrorCtx(ctx, err)
		return
	}

	if err := h.DB.Delete(&find_person).Error; err != nil {
		utils.CheckErrorCtx(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, framework.H{
		"result": ctx.Params["id"] + " deleted",
	})
}
