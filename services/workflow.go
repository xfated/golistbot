package services

import (
	"log"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func HandleUserInput(update *tgbotapi.Update) {
	/* Debugging */
	utils.LogMessage(update)
	utils.LogUpdate(update)
	utils.LogCallbackQuery(update)

	/* Check for main commands */
	message, _, err := utils.GetMessage(update)
	if err == nil {
		switch message {
		case "/start",
			"/start@toGoListBot",
			"/reset",
			"/reset@toGoListBot":
			utils.RemoveMarkupKeyboard(update, "I am ready!")
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			return
		case "/start addItem":
			// Add item in pm after redirect
			targetChat, err := utils.GetChatTarget(update)
			if err != nil {
				log.Printf("error GetChatTarget: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			if targetChat == 0 {
				utils.SendMessage(update, "Please send /additem back in the chat if you'd like to add a item")
				return
			}
			utils.SendMessage(update, "Please enter the name of the item to begin")
			if err := utils.SetUserState(update, constants.AddNewSetName); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			return
		case "/additem",
			"/additem@toGoListBot":
			// Check if is already private.
			chatID, userID, err := utils.GetChatUserID(update)
			utils.SetChatTarget(update, chatID)

			if err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			// Same == same chat
			if chatID == int64(userID) {
				utils.SendMessage(update, "Please enter the name of the item to begin")
				if err := utils.SetUserState(update, constants.AddNewSetName); err != nil {
					log.Printf("error setting state: %+v", err)
					utils.SendMessage(update, "Sorry an error occured!")
					return
				}
				return
			}
			// If not private, redirect
			utils.RedirectToBotChat(update, "Click the button to start adding", "Add item", "https://t.me/toGoListBot?start=addItem")
			return
		case "/query",
			"/query@toGoListBot":
			utils.ResetQuery(update)
			// End query if no item
			err := checkAnyItem(update)
			if err != nil {
				return
			}
			// Record id for selective force reply
			_, messageID, err := utils.GetMessage(update)
			if err != nil {
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			utils.SetMessageTarget(update, messageID)

			sendQuerySelectType(update, "What kind of query do you seek?")
			if err := utils.SetUserState(update, constants.QuerySelectType); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			return
		case "/deleteitem",
			"/deleteitem@toGoListBot":
			sendItemsToDeleteResponse(update, "Just select item do you want to delete?")
			if err := utils.SetUserState(update, constants.DeleteSelect); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			return
		case "/start editItem":
			sendItemsToEditResponse(update, "Which item would you like to edit?")
			if err := utils.SetUserState(update, constants.GetItemToEdit); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			if err := utils.SetUserState(update, constants.GetItemToEdit); err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			return
		case "/edititem",
			"/edititem@toGoListBot":
			// Check if is already private.
			chatID, userID, err := utils.GetChatUserID(update)
			utils.SetChatTarget(update, chatID)

			if err != nil {
				log.Printf("error setting state: %+v", err)
				utils.SendMessage(update, "Sorry an error occured!")
				return
			}
			// Same == same chat
			if chatID == int64(userID) {
				sendItemsToEditResponse(update, "Which item would you like to edit?")
				if err := utils.SetUserState(update, constants.GetItemToEdit); err != nil {
					log.Printf("error setting state: %+v", err)
					utils.SendMessage(update, "Sorry an error occured!")
					return
				}
				return
			}
			// If not private, redirect
			utils.RedirectToBotChat(update, "Click the button to start editing", "Edit item", "https://t.me/toGoListBot?start=editItem")
			return
		case "/help",
			"/help@toGoListBot":
			helpHandler(update)
		}
	}

	/* Get user state for Targeted handling */
	userState, err := utils.GetUserState(update)
	if err != nil {
		log.Printf("error getting user state: %+v", err)
		return
	}

	/* Idle state */
	if userState == constants.Idle {
		idleHandler(update)
		return
	}

	/* Adding new item */
	if constants.IsAddingNewItem(userState) {
		addItemHandler(update, userState)
		return
	}

	/* Querying items */
	if constants.IsQuery(userState) {
		queryHandler(update, userState)
		return
	}

	/* Delete item */
	if constants.IsDeleteItem(userState) {
		deleteItemHandler(update, userState)
		return
	}

	/* Edit item */
	if constants.IsEditItem(userState) {
		editItemHandler(update, userState)
		return
	}
}
