package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xfated/golistbot/services/constants"
	"github.com/xfated/golistbot/services/utils"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

/* Create and send template reply keyboard */
func sendTemplateReplies(update *tgbotapi.Update, text string) {
	// Create buttons
	setAddressButton := tgbotapi.NewKeyboardButton("/setAddress")
	setNotesButton := tgbotapi.NewKeyboardButton("/setNotes")
	setURLButton := tgbotapi.NewKeyboardButton("/setURL")
	addImageButton := tgbotapi.NewKeyboardButton("/addImage")
	addTagButton := tgbotapi.NewKeyboardButton("/addTag")
	removeTagButton := tgbotapi.NewKeyboardButton("/removeTag")
	previewButton := tgbotapi.NewKeyboardButton("/preview")
	submitButton := tgbotapi.NewKeyboardButton("/submit")
	cancelButton := tgbotapi.NewKeyboardButton("/cancel")
	// Create rows
	row1 := tgbotapi.NewKeyboardButtonRow(setAddressButton, setURLButton, setNotesButton)
	row2 := tgbotapi.NewKeyboardButtonRow(addImageButton, addTagButton, removeTagButton)
	row3 := tgbotapi.NewKeyboardButtonRow(cancelButton, previewButton, submitButton)

	replyKeyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	replyKeyboard.Selective = false
	utils.SetReplyMarkupKeyboard(update, text, replyKeyboard, true)
}

func sendExistingTagsResponse(update *tgbotapi.Update, text string) {
	chatID, err := utils.GetChatTarget(update)
	if err != nil {
		log.Printf("Error GetChatTarget: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!", false)
		return
	}
	chatIDString := strconv.FormatInt(chatID, 10)

	tagsMap, err := utils.GetTags(update, chatIDString)
	if err != nil {
		log.Printf("error GetTags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!", false)
		return
	}

	/* No tags, just send done */
	if len(tagsMap) == 0 {
		utils.CreateAndSendInlineKeyboard(update, "No tags found. Just click this button when you're done!", 1, "/done")
		return
	}

	/* Get already added tags */
	curTempTags, err := utils.GetTempItemTags(update)
	if err != nil {
		log.Printf("error GetTags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!", false)
		return
	}
	tags := make([]string, 0)
	i := 0
	for tag := range tagsMap {
		// if not inside current temp tags
		if !curTempTags[tag] {
			tags = append(tags, tag)
			i++
		}
	}

	tags = append(tags, "/done")
	utils.CreateAndSendInlineKeyboard(update, text, 1, tags...)
}

func sendAddedTagsResponse(update *tgbotapi.Update, text string) {
	tagsMap, err := utils.GetTempItemTags(update)
	if err != nil {
		log.Printf("error GetTags: %+v", err)
		utils.SendMessage(update, "Sorry, an error occured!", false)
	}

	/* No tags, just send done */
	if len(tagsMap) == 0 {
		utils.CreateAndSendInlineKeyboard(update, "No tags found. Just help me click that done button thanks", 1, "/done")
		return
	}

	tags := make([]string, len(tagsMap)+1)
	i := 0
	for tag := range tagsMap {
		tags[i] = tag
		i++
	}
	tags[len(tagsMap)] = "/done"
	utils.CreateAndSendInlineKeyboard(update, text, 1, tags...)
}

// func sendDoneResponse(update *tgbotapi.Update, text string) {
// 	utils.CreateAndSendInlineKeyboard(update, text, 1, "/done", "/done")
// }

func sendConfirmSubmitResponse(update *tgbotapi.Update, text string) {
	utils.CreateAndSendInlineKeyboard(update, text, 2, "yes", "no")
	// // Create buttons
	// yesButton := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	// noButton := tgbotapi.NewInlineKeyboardButtonData("no", "no")
	// // Create rows
	// row := tgbotapi.NewInlineKeyboardRow(yesButton, noButton)

	// inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	// utils.SendInlineKeyboard(update, text, inlineKeyboard)
}

func addItemHandler(update *tgbotapi.Update, userState constants.State) {
	switch userState {
	case constants.AddNewSetName:
		// Expect user to send a text message (name of item)
		// Check for slash (affect firebase query)
		if err := utils.CheckForSlash(update); err != nil {
			return
		}

		if err := utils.InitItem(update); err != nil {
			log.Printf("Error creating new item: %+v", err)
			utils.SendMessage(update, "Message should be a text", false)
			break
		}
		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			break
		}
		utils.SendMessage(update, "You may start adding the details for the item", false)
	case constants.ReadyForNextAction:
		// Expect user to select reply markup (pick next action)
		message, _, err := utils.GetMessage(update)
		if err != nil {
			log.Printf("error setting message: %+v", err)
		}
		switch message {
		case "/setAddress":
			if err := utils.SetUserState(update, constants.AddNewSetAddress); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			utils.RemoveMarkupKeyboard(update, "Send an address to be added", false)
		case "/setNotes":
			if err := utils.SetUserState(update, constants.AddNewSetNotes); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			utils.RemoveMarkupKeyboard(update, "Give some additional details as notes", false)
		case "/setURL":
			if err := utils.SetUserState(update, constants.AddNewSetURL); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			utils.RemoveMarkupKeyboard(update, "Send a URL to be added", false)
		case "/addImage":
			if err := utils.SetUserState(update, constants.AddNewSetImages); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			utils.RemoveMarkupKeyboard(update, "Send an image to be added", false)
		case "/addTag":
			if err := utils.SetUserState(update, constants.AddNewSetTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			/* Get message ID for targeted reply afterward */
			_, messageID, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			utils.SetMessageTarget(update, messageID)

			utils.RemoveMarkupKeyboard(update, "Send a tag to be added. (Can be used to query your record of items)\n"+
				"Type new or pick from existing\n\nPress \"/done\" once done!", false)
			sendExistingTagsResponse(update, "Existing tags:")
		case "/removeTag":
			if err := utils.SetUserState(update, constants.AddNewRemoveTags); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			/* Get message ID for targeted reply afterward */
			_, messageID, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			utils.SetMessageTarget(update, messageID)

			utils.RemoveMarkupKeyboard(update, "Select a tag to remove\n\nPress \"/done\" once done!", false)
			sendAddedTagsResponse(update, "Existing tags:")
		case "/preview":
			itemData, err := utils.GetTempItem(update)
			if err != nil {
				log.Printf("error getting temp item: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			utils.SendItemDetails(update, itemData, true)
			sendTemplateReplies(update, "Select your next action")
		case "/submit":
			_, messageID, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			if err := utils.SetUserState(update, constants.ConfirmAddItemSubmit); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			utils.SetMessageTarget(update, messageID)
			sendConfirmSubmitResponse(update, "Are you really ready to submit?")
		case "/cancel":
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				break
			}
			utils.RemoveMarkupKeyboard(update, "/additem process cancelled", false)
		default:
			sendTemplateReplies(update, "Please select a response from the provided options")
		}
		return
	case constants.AddNewSetAddress:
		// Expect user to send a text message (address of item)
		if err := utils.SetTempItemAddress(update); err != nil {
			log.Printf("Error adding address: %+v", err)
			utils.SendMessage(update, "Address should be a text", false)
			return
		}
		utils.SendMessage(update, fmt.Sprintf("Address set to: %s", update.Message.Text), false)

		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
	case constants.AddNewSetNotes:
		// Expect user to send a text message (notes for the item)
		if err := utils.SetTempItemNotes(update); err != nil {
			log.Printf("Error adding notes: %+v", err)
			utils.SendMessage(update, "Notes should be a text", false)
			return
		}
		utils.SendMessage(update, fmt.Sprintf("Notes set to: %s", update.Message.Text), false)

		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
	case constants.AddNewSetURL:
		// Expect user to send a text message (URL for the item)
		if err := utils.SetTempItemURL(update); err != nil {
			log.Printf("Error adding url: %+v", err)
			utils.SendMessage(update, "URL should be a text", false)
			return
		}
		utils.SendMessage(update, fmt.Sprintf("URL set to: %s", update.Message.Text), false)

		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
	case constants.AddNewSetImages:
		// Expect user to send a photo
		// should be an image input
		if err := utils.AddTempItemImage(update); err != nil {
			log.Printf("Error adding image: %+v", err)
			utils.SendMessage(update, "Error occured. Did you send an image? Try it again", false)
			return
		}
		utils.SendMessage(update, "Image added", false)

		if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
			log.Printf("error SetUserState: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
	case constants.AddNewSetTags:
		// Expect user to send a text message or Select from inline keyboard markup (set as tag for the item)
		// Check for slash (affect firebase query)
		if err := utils.CheckForSlash(update); err != nil {
			return
		}

		if update.Message != nil {
			tag, _, err := utils.GetMessage(update)
			if err != nil {
				log.Printf("error GetMessage: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			switch tag {
			case "/done",
				"done",
				"Done":
				if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
					log.Printf("error SetUserState: %+v", err)
					utils.SendMessage(update, "Sorry, an error occured!", false)
					return
				}
				// Only continue if /done is pressed
			default:
				if err := utils.AddTempItemTag(update, tag); err != nil {
					log.Printf("Error adding tag: %+v", err)
					utils.SendMessage(update, "Tag should be a text", false)
					return
				}
				utils.SendMessage(update, fmt.Sprintf("Tag \"%s\" added", tag), false)
				return
			}
		} else {
			// Then check if its a keyboard reply
			tag, err := utils.GetCallbackQueryMessage(update)
			if err != nil {
				log.Printf("error GetCallbackQueryMessage: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			if len(tag) > 0 {
				switch tag {
				case "/done":
					if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
						log.Printf("error SetUserState: %+v", err)
						utils.SendMessage(update, "Sorry, an error occured!", false)
						return
					}
				default:
					if err := utils.AddTempItemTag(update, tag); err != nil {
						log.Printf("Error adding tag: %+v", err)
						utils.SendMessage(update, "Sorry, an error occured!", false)
						return
					}
					utils.SendMessage(update, fmt.Sprintf("Tag \"%s\" added", tag), false)
					// Don't continue to next action if adding tag through inline
					return
				}
			}
		}
	case constants.AddNewRemoveTags:
		// Expect user to select from inline keyboard markup (set as tag for the item)
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options", false)
			return
		}

		tag, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}

		switch tag {
		case "/done":
			if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
		default:
			// remove tag
			if err := utils.DeleteTempItemTag(update, tag); err != nil {
				log.Printf("Error adding tag: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			utils.SendMessage(update, fmt.Sprintf("Tag \"%s\" removed", tag), false)
			sendAddedTagsResponse(update, "Existing tags:")
			// Don't continue to next action if removing tag through inline
			return
		}
	case constants.ConfirmAddItemSubmit:
		// Expect user to select from inline query (yes or no to submit)
		/* If user send a message instead */
		if update.Message != nil {
			utils.SendMessage(update, "Please select from the above options", false)
			return
		}

		confirm, err := utils.GetCallbackQueryMessage(update)
		if err != nil {
			log.Printf("error getting message from callback: %+v", err)
			utils.SendMessage(update, "Sorry, an error occured!", false)
			return
		}
		if confirm == "yes" {
			// Get target chat, where additem was initiated
			chatID, err := utils.GetChatTarget(update)
			chatIDString := strconv.FormatInt(chatID, 10)
			if err != nil {
				log.Printf("error getting message from callback: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}

			// Submit
			name, err := utils.AddItemFromTemp(update, chatIDString)
			if err != nil {
				log.Printf("error adding item from temp: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			if err := utils.SetUserState(update, constants.Idle); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			utils.RemoveMarkupKeyboard(update, fmt.Sprintf("%s has been added/edited!", name), false)
			utils.SendMessage(update, "To add/edit a new item to any chat, please initiate /additem or /edititem in that chat", false)
			err = utils.SetChatTarget(update, 0)
			if err != nil {
				log.Printf("error SetChatTarget: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
			return
		} else if confirm == "no" {
			if err := utils.SetUserState(update, constants.ReadyForNextAction); err != nil {
				log.Printf("error SetUserState: %+v", err)
				utils.SendMessage(update, "Sorry, an error occured!", false)
				return
			}
		}
	}

	/* Create and send keyboard for targeted response */
	sendTemplateReplies(update, "What do you want to do next?")
}
